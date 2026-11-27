import GoLean.Math.Big.IntBitwise
import GoDiff.Math.Bits.Tuple

open GoLean.Math.Big

private def render [ToString α] : Except BitFault α → String
  | .ok x => toString x
  | .error .negativeIndex => "negative bit index"
  | .error .invalidBit => "set bit is not 0 or 1"

def main : IO Unit := runTuple fun fields => do
  let [x, y, shift, index, bit] := fields | throw (IO.userError "expected integer operands, shift, index and bit")
  return String.intercalate "|" [toString (intAnd x y), toString (intAndNot x y),
    toString (intOr x y), toString (intXor x y), toString (intNot x), toString (intLsh x shift.toNat),
    toString (intRsh x shift.toNat), render (intBit x index), render (intSetBit x index bit.toNat),
    toString (intTrailingZeroBits x)]
