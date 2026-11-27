import GoLean.Bytes.SplitAfter
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let _ :: size :: spare :: raw := row | throw (IO.userError "expected limit, separator size, spare capacity and bytes")
  let words ← raw.mapM (wordInput 8)
  let bytes := words.map BitVec.toNat
  match GoLean.Bytes.splitAfter (bytes.drop size.toNat) (bytes.take size.toNat) spare.toNat with
  | none => return "nil"
  | some parts =>
    return "slice" ++ String.intercalate ";" (parts.map fun (part, capacity) =>
      "f" ++ toString capacity ++ ":" ++ String.intercalate "," (part.map toString))
