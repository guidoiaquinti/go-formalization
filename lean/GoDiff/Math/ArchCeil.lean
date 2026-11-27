import GoLean.Math.ArchCeil
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, raw] := fields | throw (IO.userError "expected a profile and binary64 word")
  let word := GoLean.Math.archCeil (profile = 1) (← wordInput 64 raw)
  return if word % 2 ^ 63 > 0x7ff0000000000000 then "NaN" else toString word
