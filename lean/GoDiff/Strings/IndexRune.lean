import GoLean.Strings.IndexRune
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let r :: words := fields | throw (IO.userError "expected rune and input bytes")
  let input ← words.mapM fun n => do
    return (← wordInput 8 n).toNat
  let offset := (GoLean.Strings.indexRune input r).map toString |>.getD "-1"
  return s!"{offset} {GoLean.Strings.containsRune input r}"
