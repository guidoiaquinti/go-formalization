import GoLean.Strings.Count
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let length :: words := fields | throw (IO.userError "expected separator length")
  if length < 0 || words.length < length.toNat then
    throw (IO.userError "invalid separator length")
  let bytes ← words.mapM (fun n => do return (← wordInput 8 n).toNat)
  return toString (GoLean.Strings.count (bytes.drop length.toNat) (bytes.take length.toNat))
