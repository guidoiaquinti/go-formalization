import GoLean.Image.Color.SqDiff
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [x, y] := fields | throw (IO.userError "expected two uint32 words")
  let x ← wordInput 32 x
  let y ← wordInput 32 y
  return toString (GoLean.Image.Color.sqDiff x y).toNat
