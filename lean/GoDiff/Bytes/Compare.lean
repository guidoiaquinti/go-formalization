import GoLean.Bytes.Compare
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

def main : IO Unit := runTuple fun fields => do
  let leftLen :: rightLen :: words := fields | throw (IO.userError "expected two slice lengths")
  if leftLen < -1 || rightLen < -1 || leftLen.toNat + rightLen.toNat != words.length then
    throw (IO.userError "invalid slice lengths")
  let data ← words.mapM (wordInput 8)
  let left : Slice := if leftLen == -1 then none else some (data.take leftLen.toNat, leftLen.toNat)
  let right : Slice := if rightLen == -1 then none else some (data.drop leftLen.toNat, rightLen.toNat)
  return s!"{equal left right} {GoLean.Bytes.compare left right}"
