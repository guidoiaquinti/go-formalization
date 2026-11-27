import GoLean.Internal.Strconv.Bool
import GoDiff.Math.Bits.Tuple

open GoLean.Internal.Strconv

private def wire (raw : List UInt8) :=
  String.intercalate "," (raw.map (toString ·.toNat))

def main : IO Unit := runTuple fun row => do
  let b :: tag :: length :: capacity :: words := row |
    throw (IO.userError "expected boolean, error tag and slice geometry")
  let words ← words.mapM fun n => return UInt8.ofNat (← wordInput 8 n).toNat
  let (value, bad) := parseBool (words.drop capacity.toNat)
  let appended := appendBool (words.take capacity.toNat) length.toNat (b == 1)
  return s!"{value}|{bad}|{formatBool (b == 1)}|{wire appended.result}|{wire appended.oldBacking}|{appended.reused}|{errorText tag}"
