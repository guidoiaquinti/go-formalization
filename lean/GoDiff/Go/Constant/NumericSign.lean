import GoLean.Go.Constant.NumericSign
import GoDiff.Math.Bits.Tuple

open GoLean.Go.Constant

private def report : Except String Int → String
  | .ok sign => toString sign
  | .error message => "panic:" ++ message

def main : IO Unit := runTuple fun fields => do
  let [x, y, mode] := fields | throw (IO.userError "expected signed numeric components and mode")
  let xi := SignValue.int64 x
  let yi := SignValue.rational y
  let complex := SignValue.complex xi yi
  let alternate := if mode == 0 then SignValue.unknown
    else if mode == 1 then .other "<nil>" else .other "true"
  return String.intercalate "|" [
    toString (int64ValKind x), toString (intValKind x), toString (ratValKind x),
    toString (floatValKind x), toString (complexValKind xi yi),
    report (constantSign xi), report (constantSign (.bigInt x)),
    report (constantSign (.rational x)), report (constantSign (.float x)),
    report (constantSign complex), report (constantSign (.complex complex (.float (-x)))),
    report (constantSign alternate), report (constantSign (.complex alternate yi))]
