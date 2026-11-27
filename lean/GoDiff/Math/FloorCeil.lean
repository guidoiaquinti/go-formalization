import GoLean.Math.FloorCeil
import GoDiff.Math.Bits.Tuple

private def observe (word : Nat) : String :=
  if word % 2 ^ 63 > 0x7ff0000000000000 then "NaN" else toString word

def main : IO Unit := runTuple fun fields => do
  let [raw] := fields | throw (IO.userError "expected one binary64 word")
  let word ← wordInput 64 raw
  return observe (GoLean.Math.floorWord word) ++ " " ++ observe (GoLean.Math.ceilWord word)
