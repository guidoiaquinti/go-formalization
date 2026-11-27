import GoLean.Internal.Bytealg.Compare
import GoDiff.Math.Bits.Tuple

open GoLean.Internal.Bytealg

def main : IO Unit := runTuple fun fields => do
  let leftLen :: rightLen :: words := fields | throw (IO.userError "expected two slice lengths")
  if leftLen < -1 || rightLen < -1 || leftLen.toNat + rightLen.toNat != words.length then
    throw (IO.userError "invalid slice lengths")
  let data ← words.mapM (wordInput 8)
  let left := data.take leftLen.toNat
  let right := data.drop leftLen.toNat
  let a : GoLean.Bytes.Slice := if leftLen == -1 then none else some (left, left.length)
  let b : GoLean.Bytes.Slice := if rightLen == -1 then none else some (right, right.length)
  return s!"{compare a b},{compareString left right},{runtimeCmpstring left right},{equal a b}"
