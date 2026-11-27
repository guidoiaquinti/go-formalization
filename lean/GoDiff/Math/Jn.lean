import GoLean.Math.Jn
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, order, raw] := fields |
    throw (IO.userError "expected profile, order and binary64 word")
  if 64 < order.natAbs then throw (IO.userError "order outside modeled bounds")
  let word ← wordInput 64 raw
  match GoLean.Math.jnWord (profile == 2) (profile == 1) order word with
  | none => return "fuel"
  | some none => return "NaN"
  | some (some result) => return toString result
