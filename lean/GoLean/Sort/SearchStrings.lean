import GoLean.Sort.SearchInterval
import Init.Data.List.Lex

namespace GoLean.Sort

/-- Go 1.27.0 src/sort/search.go:135-141. Compare Go strings as unsigned byte lists. -/
def searchStrings (xs : List (List UInt8)) (key : List UInt8) : Nat :=
  searchInterval (fun i => !((xs[i]?.getD []).lex key)) 0 xs.length

theorem searchStrings_spec (xs : List (List UInt8)) (key : List UInt8) (cut : Nat)
    (bounded : cut ≤ xs.length)
    (partition : ∀ i < xs.length, (¬xs[i]?.getD [] < key) ↔ cut ≤ i) :
    searchStrings xs key = cut := by
  apply searchInterval_partition _ 0 xs.length cut (Nat.zero_le _) bounded
  intro i _ hi
  simpa [Bool.eq_false_iff] using partition i hi


end GoLean.Sort
