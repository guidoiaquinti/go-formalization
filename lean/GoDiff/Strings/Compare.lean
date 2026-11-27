import GoLean.Strings.Compare
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let left :: right :: words := fields | throw (IO.userError "expected two byte counts")
  if left < -1 || right < -1 || left.toNat + right.toNat != words.length then
    throw (IO.userError "invalid byte counts")
  let bytes ← words.mapM (wordInput 8)
  return toString (GoLean.Strings.compare (bytes.take left.toNat) (bytes.drop left.toNat))
