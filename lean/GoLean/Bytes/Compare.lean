import GoLean.Bytes.Affix
import GoLean.Slices.Compare

namespace GoLean.Bytes

/-- Go compares the string values of both slices, treating nil as empty.
    https://github.com/golang/go/blob/go1.27.0/src/bytes/bytes.go#L17-L23 -/
def equal (a b : Slice) : Bool := decide (contents a = contents b)

/-- Compare unsigned byte values in order, stopping at the first difference
    or at the end of the shorter slice. The model follows the generic scan.
    https://github.com/golang/go/blob/go1.27.0/src/bytes/bytes.go#L25-L30
    https://github.com/golang/go/blob/go1.27.0/src/internal/bytealg/compare_generic.go#L11-L36 -/
def compare (a b : Slice) : Int :=
  let unsigned (s : Slice) := (contents s).map fun x => (x.toNat : Int)
  Slices.sliceCompare Slices.compareInt (unsigned a) (unsigned b)

theorem equal_spec (a b : Slice) :
    equal a b = true ↔ contents a = contents b := by
  simp [equal]

theorem compare_spec (a b : Slice) :
    let x := (contents a).map fun v => (v.toNat : Int)
    let y := (contents b).map fun v => (v.toNat : Int)
    (compare a b = -1 ↔ x < y) ∧
    (compare a b = 0 ↔ x = y) ∧
    (compare a b = 1 ↔ y < x) :=
  Slices.sliceCompare_int_spec _ _

end GoLean.Bytes
