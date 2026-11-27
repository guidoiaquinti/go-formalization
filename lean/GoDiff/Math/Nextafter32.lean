import GoLean.Math.Nextafter32
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [x, y] := fields | throw (IO.userError "expected two binary32 words")
  let a ← wordInput 32 x
  let b ← wordInput 32 y
  return toString (GoLean.Math.nextafter32 a.toNat b.toNat)
