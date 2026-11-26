import GoLean.Sort.IndexedSwap

namespace GoLean.Sort

-- Each UInt64 represents the exact bits of a binary64 element.

/-- Go 1.27.0 src/sort/sort.go:145 checks destinations in index order before either write.
A failed bounds check leaves the input list untouched. -/
def float64SliceSwap (xs : List UInt64) (i j : Int) : Except Int (List UInt64) := indexedSwap xs i j

theorem float64SliceSwap_spec (xs : List UInt64) (i j : Int) :
    float64SliceSwap xs i j =
      if i < 0 ∨ xs.length ≤ i.toNat then .error i
      else if j < 0 ∨ xs.length ≤ j.toNat then .error j
      else .ok ((xs.set i.toNat (xs[j.toNat]?.getD 0)).set j.toNat (xs[i.toNat]?.getD 0)) := indexedSwap_spec 0 xs i j

theorem float64SliceSwap_index_spec (xs : List UInt64) (i j : Int)
    (hi : 0 ≤ i) (hj : 0 ≤ j) (hib : i.toNat < xs.length) (hjb : j.toNat < xs.length)
    (output : List UInt64) (h : float64SliceSwap xs i j = .ok output) (k : Nat) :
    output.length = xs.length ∧
    output[k]? = xs[if k = j.toNat then i.toNat else if k = i.toNat then j.toNat else k]? :=
  indexedSwap_index_spec 0 xs i j hi hj hib hjb output h k

end GoLean.Sort
