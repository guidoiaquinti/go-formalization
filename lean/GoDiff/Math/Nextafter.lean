import GoLean.Math.Nextafter
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [x, y] := fields | throw (IO.userError "expected two binary64 words")
  let a ← wordInput 64 x
  let b ← wordInput 64 y
  return toString (GoLean.Math.nextafter a.toNat b.toNat)
