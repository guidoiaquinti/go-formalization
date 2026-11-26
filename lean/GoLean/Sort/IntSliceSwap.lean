import GoLean.Sort.IndexedSwap

namespace GoLean.Sort

/-- Go 1.27.0 src/sort/sort.go:127 checks destinations in index order before either write.
A failed bounds check leaves the input list untouched. -/
def intSliceSwap (xs : List Int) (i j : Int) : Except Int (List Int) := indexedSwap xs i j

theorem intSliceSwap_spec (xs : List Int) (i j : Int) :
    intSliceSwap xs i j =
      if i < 0 ∨ xs.length ≤ i.toNat then .error i
      else if j < 0 ∨ xs.length ≤ j.toNat then .error j
      else .ok ((xs.set i.toNat (xs[j.toNat]?.getD 0)).set j.toNat (xs[i.toNat]?.getD 0)) := indexedSwap_spec 0 xs i j

theorem intSliceSwap_index_spec (xs : List Int) (i j : Int)
    (hi : 0 ≤ i) (hj : 0 ≤ j) (hib : i.toNat < xs.length) (hjb : j.toNat < xs.length)
    (output : List Int) (h : intSliceSwap xs i j = .ok output) (k : Nat) :
    output.length = xs.length ∧
    output[k]? = xs[if k = j.toNat then i.toNat else if k = i.toNat then j.toNat else k]? :=
  indexedSwap_index_spec 0 xs i j hi hj hib hjb output h k

end GoLean.Sort
