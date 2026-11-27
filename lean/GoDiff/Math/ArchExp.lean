import GoLean.Math.ArchExp
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, word] := fields | throw (IO.userError "expected profile and binary64 word")
  let input ← wordInput 64 word
  return toString (GoLean.Math.archExpWord (profile == 2) (profile == 1) input)
