import Init.Data.List.Lemmas

namespace GoLean.Internal.Dag

private abbrev Edges := List (Nat × Nat)

private def reduceVia (nodes : List Nat) (edges : Edges) (source via : Nat) : Edges :=
  if edges.contains (source, via) then
    nodes.foldl (fun current dest =>
      if current.contains (via, dest) then
        current.filter (· != (source, dest))
      else current) edges
  else edges

/-- The source scans the nodes in insertion order. Lookups use the current edge
    set, and each deletion removes the selected source–destination pair.
    https://github.com/golang/go/blob/go1.27.0/src/internal/dag/alg.go#L48-L63
    https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L75-L81 -/
def transitiveReduction (nodes : List Nat) (edges : List (Nat × Nat)) : List (Nat × Nat) :=
  nodes.foldl (fun current source => nodes.foldl (fun acc via =>
    reduceVia nodes acc source via) current) edges

private def preservesCovers (original current : Edges) : Prop :=
  current ⊆ original ∧
  ∀ source dest, (source, dest) ∈ original →
    (∀ via, (source, via) ∈ original → (via, dest) ∉ original) →
    (source, dest) ∈ current

private theorem fold_preserves {α β : Type} (items : List α) (state : β)
    (step : β → α → β) (p : β → Prop) (initial : p state)
    (preserved : ∀ s a, p s → p (step s a)) : p (items.foldl step state) := by
  induction items generalizing state with
  | nil => exact initial
  | cons head tail ih => exact ih (step state head) (preserved state head initial)

private theorem reduceVia_preserves (nodes : List Nat) (original current : Edges)
    (source via : Nat) (invariant : preservesCovers original current) :
    preservesCovers original (reduceVia nodes current source via) := by
  unfold reduceVia
  split
  next sourceEdge =>
    have sourcePresent : (source, via) ∈ original :=
      invariant.1 (List.contains_iff_mem.mp (by simpa using sourceEdge))
    apply fold_preserves _ _ _ _ invariant
    intro state target ih
    split
    next targetEdge =>
      have targetPresent : (via, target) ∈ original :=
        ih.1 (List.contains_iff_mem.mp (by simpa using targetEdge))
      constructor
      · intro edge added
        exact ih.1 (List.mem_filter.mp added).1
      · intro src dest present cover
        apply List.mem_filter.mpr
        refine ⟨ih.2 src dest present cover, ?_⟩
        apply bne_iff_ne.mpr
        intro eq
        rcases Prod.mk.inj eq with ⟨rfl, rfl⟩
        exact cover via sourcePresent targetPresent
    next _ => exact ih
  next _ => exact invariant

/-- The scan only deletes edges. An original edge without an original two-hop
    witness is retained. -/
theorem transitiveReduction_spec (nodes : List Nat) (edges : List (Nat × Nat)) :
    transitiveReduction nodes edges ⊆ edges ∧
    ∀ source dest, (source, dest) ∈ edges →
      (∀ via, (source, via) ∈ edges → (via, dest) ∉ edges) →
      (source, dest) ∈ transitiveReduction nodes edges := by
  change preservesCovers edges (transitiveReduction nodes edges)
  unfold transitiveReduction
  apply fold_preserves
  · exact ⟨fun _ h => h, fun _ _ h _ => h⟩
  · intro current source ih
    apply fold_preserves _ _ _ _ ih
    intro acc via stable
    exact reduceVia_preserves nodes edges acc source via stable

end GoLean.Internal.Dag
