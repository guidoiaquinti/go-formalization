import GoLean.Path.Filepath.SameWord
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

def main : IO Unit := runTuple fun fields => do
  let leftLength :: raw := fields | throw (IO.userError "expected left length and comparison bytes")
  if leftLength < 0 || leftLength.toNat > raw.length then
    throw (IO.userError "invalid left length")
  let bytes ← raw.mapM (wordInput 8)
  let input := bytes.map BitVec.toNat
  let left := input.take leftLength.toNat
  let right := input.drop leftLength.toNat
  return s!"{sameWord false left right}|{sameWord true left right}"
