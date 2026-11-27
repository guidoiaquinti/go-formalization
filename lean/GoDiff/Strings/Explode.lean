import GoLean.Strings.Explode
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let limit :: raw := row | throw (IO.userError "expected limit and bytes")
  let bytes ← raw.mapM (wordInput 8)
  let parts := GoLean.Strings.explodeDirect (bytes.map BitVec.toNat) limit
  return "slice" ++ String.intercalate ";" (parts.map fun part =>
    "f" ++ String.intercalate "," (part.map toString))
