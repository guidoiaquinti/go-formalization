import GoLean.Slices.Compare
import Init.Data.BitVec.Basic

namespace GoLean.Strings

/-- The generic Go comparison scans unsigned bytes to the first difference,
    then orders equal prefixes by length.
    https://github.com/golang/go/blob/go1.27.0/src/strings/compare.go#L8-L17
    https://github.com/golang/go/blob/go1.27.0/src/internal/bytealg/compare_generic.go#L38-L40
    https://github.com/golang/go/blob/go1.27.0/src/internal/bytealg/compare_generic.go#L54-L76 -/
def compare (a b : List (BitVec 8)) : Int :=
  let unsigned (bytes : List (BitVec 8)) := bytes.map fun v => (v.toNat : Int)
  Slices.sliceCompare Slices.compareInt (unsigned a) (unsigned b)

theorem compare_spec (a b : List (BitVec 8)) :
    let x := a.map fun v => (v.toNat : Int)
    let y := b.map fun v => (v.toNat : Int)
    (compare a b = -1 ↔ x < y) ∧
    (compare a b = 0 ↔ x = y) ∧
    (compare a b = 1 ↔ y < x) :=
  Slices.sliceCompare_int_spec _ _

end GoLean.Strings
