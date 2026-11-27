import GoLean.Internal.Bytealg.RuntimeMemory
import GoDiff.Math.Bits.Tuple

open GoLean.Internal.Bytealg

def main : IO Unit := runTuple fun fields => do
  let length :: size :: leftLength :: raw := fields |
    throw (IO.userError "expected allocation length, comparison size and left length")
  let data ← raw.mapM (wordInput 8)
  let left := data.take leftLength.toNat
  let right := data.drop leftLength.toNat
  let allocation := match makeNoZero length 281474976710656 with
    | none => "panic"
    | some result => s!"{result.length}:{result.capacity}:{result.nonNil}"
  return s!"{allocation}|{runtimeMemequal size.toNat left right}"
