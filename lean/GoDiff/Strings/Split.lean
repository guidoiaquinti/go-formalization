import GoLean.Strings.Split
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let size :: raw := row | throw (IO.userError "expected separator size and bytes")
  let words ← raw.mapM (wordInput 8)
  let bytes := words.map BitVec.toNat
  match GoLean.Strings.split (bytes.drop size.toNat) (bytes.take size.toNat) with
  | none => return "nil"
  | some parts =>
    return "slice" ++ String.intercalate ";" (parts.map fun part =>
      "f" ++ String.intercalate "," (part.map toString))
