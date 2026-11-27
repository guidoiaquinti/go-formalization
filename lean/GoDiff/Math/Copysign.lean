import GoLean.Math.Copysign
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [magnitude, sign] := fields | throw (IO.userError "expected two binary64 words")
  let magnitudeWord ← wordInput 64 magnitude
  let signWord ← wordInput 64 sign
  return toString (GoLean.Math.copysign magnitudeWord signWord).toNat
