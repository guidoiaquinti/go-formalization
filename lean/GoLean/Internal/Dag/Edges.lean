import GoLean.Internal.Dag.HasEdge
import Init.Data.List.Sort.Lemmas

namespace GoLean.Internal.Dag

/-- A missing receiver represents a panic. Otherwise Go collects every key
in the selected row, even when its value is false, and orders the keys by
the signed rank read from byLabel. Missing ranks have value zero. Merge sort
models the key order when destination ranks are distinct.
https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L83-L92 -/
def edges (receiver : Option (Option EdgeMap)) (rank : String → Int)
    (source : String) : Option (List String) := do
  let outer ← receiver
  let row := (outer.bind fun rows => rows.lookup source).join
  let keys := (row.getD []).map Prod.fst
  return keys.mergeSort (fun a b => decide (rank a ≤ rank b))

/-- The observable keys before ordering. A missing source or nil map has
    no destinations. -/
def destinationKeys (outer : Option EdgeMap) (source : String) : List String :=
  match outer with
  | none => []
  | some rows =>
    match rows.lookup source with
    | some (some row) => row.map Prod.fst
    | _ => []

theorem edges_panics_iff (receiver : Option (Option EdgeMap)) (rank : String → Int)
    (source : String) : edges receiver rank source = none ↔ receiver = none := by
  cases receiver <;> simp [edges]

private theorem edges_eq_sort (outer : Option EdgeMap) (rank : String → Int)
    (source : String) :
    edges (some outer) rank source =
      some ((destinationKeys outer source).mergeSort (fun a b => decide (rank a ≤ rank b))) := by
  cases outer with
  | none => rfl
  | some rows =>
    cases h : rows.lookup source with
    | none => simp [edges, destinationKeys, h]
    | some row => cases row <;> simp [edges, destinationKeys, h]

theorem edges_spec (outer : Option EdgeMap) (rank : String → Int)
    (source : String) (result : List String)
    (distinctRanks : ∀ a ∈ destinationKeys outer source, ∀ b ∈ destinationKeys outer source,
      rank a = rank b → a = b) :
    edges (some outer) rank source = some result ↔
      result.Perm (destinationKeys outer source) ∧
      result.Pairwise (fun a b => rank a ≤ rank b) := by
  let keys := destinationKeys outer source
  let sorted := keys.mergeSort (fun a b => decide (rank a ≤ rank b))
  have order : sorted.Pairwise (fun a b => rank a ≤ rank b) := by
    have trans (a b c : String) : decide (rank a ≤ rank b) → decide (rank b ≤ rank c) →
        decide (rank a ≤ rank c) := by simp; exact Int.le_trans
    have total (a b : String) : (decide (rank a ≤ rank b) || decide (rank b ≤ rank a)) := by
      simpa using Int.le_total (rank a) (rank b)
    simpa [sorted] using List.pairwise_mergeSort trans total keys
  have perm : sorted.Perm keys := List.mergeSort_perm _ _
  rw [edges_eq_sort]
  change some sorted = some result ↔ result.Perm keys ∧ result.Pairwise _
  constructor
  · intro h
    rw [← Option.some.inj h]
    exact ⟨perm, order⟩
  · rintro ⟨resultPerm, resultOrder⟩
    congr 1
    exact (perm.trans resultPerm.symm).eq_of_pairwise
      (fun a b ha hb hab hba =>
        distinctRanks a (perm.mem_iff.mp ha) b (resultPerm.mem_iff.mp hb) (Int.le_antisymm hab hba))
      order resultOrder

end GoLean.Internal.Dag
