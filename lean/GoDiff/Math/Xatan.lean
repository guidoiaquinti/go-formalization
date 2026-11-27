import GoLean.Math.Xatan
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, word] := fields | throw (IO.userError "expected profile and binary64 word")
  return toString (GoLean.Math.xatanWord (profile == 1) (← wordInput 64 word))
