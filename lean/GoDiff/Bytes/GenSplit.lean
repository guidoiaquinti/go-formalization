import GoLean.Bytes.GenSplit
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let limit :: save :: size :: spareCount :: raw := row |
    throw (IO.userError "expected limit, save, separator size, spare capacity and bytes")
  let words ← raw.mapM (wordInput 8)
  let bytes := words.map BitVec.toNat
  let input := bytes.drop size.toNat
  match GoLean.Bytes.genSplit input (bytes.take size.toNat) (List.replicate spareCount.toNat 165) save limit with
  | .error _ => return "panic"
  | .ok none => return "nil"
  | .ok (some parts) =>
    return "slice" ++ String.intercalate ";" (parts.map fun (p, c) =>
      "f" ++ toString c ++ ":" ++ String.intercalate "," (p.map toString))
