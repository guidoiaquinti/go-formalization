import GoLean.Math.Yn
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, order, raw] := fields |
    throw (IO.userError "expected profile, order and binary64 word")
  if 64 < order.natAbs then throw (IO.userError "order outside modeled bounds")
  let word ← wordInput 64 raw
  let result := GoLean.Math.ynWord (profile == 2) (profile == 1) order word
  return (result.map toString).getD "NaN"
