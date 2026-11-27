import GoLean.Strings.GenSplit
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let limit :: save :: size :: raw := row | throw (IO.userError "expected limit, save, separator size and bytes")
  let words ← raw.mapM (wordInput 8)
  let bytes := words.map BitVec.toNat
  match GoLean.Strings.genSplitDirect (bytes.drop size.toNat) (bytes.take size.toNat) save limit with
  | .error () => return "panic"
  | .ok none => return "nil"
  | .ok (some parts) =>
    return "slice" ++ String.intercalate ";" (parts.map fun part =>
      "f" ++ String.intercalate "," (part.map toString))
