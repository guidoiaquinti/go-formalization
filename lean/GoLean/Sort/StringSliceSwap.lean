import GoLean.Sort.IndexedSwap

namespace GoLean.Sort

/-- Go 1.27.0 src/sort/sort.go:160 checks destinations in index order before either write.
A failed bounds check leaves the input list untouched. -/
def stringSliceSwap (xs : List (List UInt8)) (i j : Int) : Except Int (List (List UInt8)) := indexedSwap xs i j

theorem stringSliceSwap_spec (xs : List (List UInt8)) (i j : Int) :
    stringSliceSwap xs i j =
      if i < 0 ∨ xs.length ≤ i.toNat then .error i
      else if j < 0 ∨ xs.length ≤ j.toNat then .error j
      else .ok ((xs.set i.toNat (xs[j.toNat]?.getD [])).set j.toNat (xs[i.toNat]?.getD [])) := indexedSwap_spec [] xs i j

theorem stringSliceSwap_index_spec (xs : List (List UInt8)) (i j : Int)
    (hi : 0 ≤ i) (hj : 0 ≤ j) (hib : i.toNat < xs.length) (hjb : j.toNat < xs.length)
    (output : List (List UInt8)) (h : stringSliceSwap xs i j = .ok output) (k : Nat) :
    output.length = xs.length ∧
    output[k]? = xs[if k = j.toNat then i.toNat else if k = i.toNat then j.toNat else k]? :=
  indexedSwap_index_spec [] xs i j hi hj hib hjb output h k

end GoLean.Sort
