import GoLean.Sort.SearchInterval

namespace GoLean.Sort

/-- Go 1.27.0 src/sort/search.go:58-72. Nonpositive lengths leave
the interval empty. The predicate describes a normally returning callback. -/
def search (n : Int) (test : Nat → Bool) : Nat :=
  searchInterval test 0 n.toNat

theorem search_spec (n : Int) (test : Nat → Bool) (cut : Nat)
    (bounded : cut ≤ n.toNat)
    (partition : ∀ i < n.toNat, test i = true ↔ cut ≤ i) :
    search n test = cut :=
  searchInterval_partition test 0 n.toNat cut (Nat.zero_le _) bounded
    (fun i _ hi => partition i hi)


end GoLean.Sort
