import GoLean.Go.Constant.NumericConversions
import GoDiff.Math.Bits.Tuple

open GoLean.Go.Constant

private def rational (r : RationalCell) : String := s!"{r.numerator}/{r.denominator}"
private def float (f : FloatCell) : String := s!"{f.integerValue}/{f.precision}"
private def ref : ConstantRef → String
  | .absent => "nil"
  | .typedNil => "typedNil"
  | .identity id => toString id

def main : IO Unit := runTuple fun fields => do
  let [integer, big, tag] := fields | throw (IO.userError "expected integer, magnitude and tag")
  let x := if tag == 0 then ConstantRef.absent else if tag == 1 then .typedNil else .identity tag.toNat
  let z := vtoc x
  return String.intercalate "|" [toString (newNumericInt 1).value, rational (newNumericRat 2),
    float (newNumericFloat 3), toString (i64toi 4 integer).value, rational (i64tor 5 integer),
    float (i64tof 6 integer), rational (itor 7 ⟨8,big⟩), s!"{ref z.real}:{z.imaginary}"]
