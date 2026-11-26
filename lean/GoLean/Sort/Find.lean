import GoLean.Sort.SearchInterval

namespace GoLean.Sort

/-- Go 1.27.0 src/sort/search.go:99-115. The final comparison executes only
when the insertion point lies inside the search range. -/
def sortFind (n : Int) (cmp : Nat → Int) : Nat × Bool :=
  let i := searchInterval (fun j => decide (cmp j ≤ 0)) 0 n.toNat
  (i, decide (Int.ofNat i < n) && (cmp i == 0))

theorem sortFind_spec (n : Int) (cmp : Nat → Int) (cut : Nat)
    (bounded : cut ≤ n.toNat)
    (partition : ∀ i < n.toNat, cmp i ≤ 0 ↔ cut ≤ i) :
    sortFind n cmp = (cut, decide (Int.ofNat cut < n ∧ cmp cut = 0)) := by
  have foundCut : searchInterval (fun i => decide (cmp i ≤ 0)) 0 n.toNat = cut := by
    apply searchInterval_partition _ 0 n.toNat cut (Nat.zero_le _) bounded
    intro i _ hi
    simpa using partition i hi
  rw [sortFind, foundCut]
  congr 1


end GoLean.Sort
