import Init.Data.List.Sort.Lemmas

namespace GoLean.Internal.Dag

/-- The accumulated list records completed marks and reverse postorder.
    Go checks the mark, visits ordered destinations, then marks and appends
    the node. Consing here combines append and the final in-place reversal.
    Fuel bounds the recursion; on the supported DAGs it is the node count.
    https://github.com/golang/go/blob/go1.27.0/src/internal/dag/alg.go#L23-L47
    https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L83-L92 -/
private def visitTopo (next : Nat → List Nat) : Nat → List Nat → Nat → List Nat
  | 0, done, _ => done
  | fuel + 1, done, node =>
    if node ∈ done then done
    else node :: ((next node).mergeSort (· ≤ ·)).foldl (visitTopo next fuel) done

/-- Node identifiers encode the distinct insertion ranks used by Edges. -/
def topo (roots : List Nat) (next : Nat → List Nat) : List Nat :=
  roots.foldl (visitTopo next roots.length) []

private theorem visitTopo_invariant (next : Nat → List Nat) (height : Nat → Nat)
    (increases : ∀ src, ∀ dest ∈ next src, height src < height dest)
    (fuel : Nat) (done : List Nat) (node : Nat) (unique : done.Nodup) :
    let result := visitTopo next fuel done node
    result.Nodup ∧ (∀ n ∈ done, n ∈ result) ∧
      (∀ n ∈ result, n ∈ done ∨ height node ≤ height n) := by
  induction fuel generalizing done node with
  | zero => exact ⟨unique, fun _ h => h, fun _ h => Or.inl h⟩
  | succ fuel ih =>
    by_cases marked : node ∈ done
    · simpa [visitTopo, marked] using
        (show done.Nodup ∧ (∀ n ∈ done, n ∈ done) ∧
          (∀ n ∈ done, n ∈ done ∨ height node ≤ height n) from
          ⟨unique, fun _ h => h, fun _ h => Or.inl h⟩)
    · let ordered := (next node).mergeSort (· ≤ ·)
      let result := ordered.foldl (visitTopo next fuel) done
      have visited : result.Nodup ∧ (∀ n ∈ done, n ∈ result) ∧
          (∀ n ∈ result, n ∈ done ∨ height node < height n) := by
        apply List.foldlRecOn (motive := fun (acc : List Nat) =>
          acc.Nodup ∧ (∀ n ∈ done, n ∈ acc) ∧
            (∀ n ∈ acc, n ∈ done ∨ height node < height n))
        · exact ⟨unique, fun _ h => h, fun _ h => Or.inl h⟩
        · intro acc h child mem
          have edge : child ∈ next node := (List.mergeSort_perm _ _).mem_iff.mp mem
          rcases ih _ _ h.1 with ⟨unique', preserves, added⟩
          refine ⟨unique', fun n hn => preserves n (h.2.1 n hn), ?_⟩
          intro n hn
          rcases added n hn with existing | higher
          · exact h.2.2 n existing
          · exact Or.inr (Nat.lt_of_lt_of_le (increases node child edge) higher)
      rw [visitTopo, ite_eq_right marked]
      change (node :: result).Nodup ∧ _ ∧ _
      refine ⟨List.nodup_cons.mpr ⟨?_, visited.1⟩, ?_, ?_⟩
      · intro h
        rcases visited.2.2 node h with old | impossible
        · exact marked old
        · exact Nat.lt_irrefl _ impossible
      · intro n h
        exact List.mem_cons_of_mem _ (visited.2.1 n h)
      · intro n h
        rcases List.mem_cons.mp h with rfl | inResult
        · exact Or.inr (Nat.le_refl _)
        · exact (visited.2.2 n inResult).imp_right Nat.le_of_lt

theorem topo_nodup (roots : List Nat) (next : Nat → List Nat) (height : Nat → Nat)
    (increases : ∀ src, ∀ dest ∈ next src, height src < height dest) :
    (topo roots next).Nodup := by
  unfold topo
  apply List.foldlRecOn (motive := List.Nodup)
  · exact List.nodup_nil
  · intro acc unique node _
    exact (visitTopo_invariant next height increases roots.length acc node unique).1

end GoLean.Internal.Dag
