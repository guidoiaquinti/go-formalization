import GoLean.Bytes.Explode
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let limit :: spare :: raw := row | throw (IO.userError "expected limit, spare capacity and bytes")
  let words ← raw.mapM (wordInput 8)
  match GoLean.Bytes.explode (words.map BitVec.toNat) spare.toNat limit with
  | none => return "nil"
  | some parts =>
    return "slice" ++ String.intercalate ";" (parts.map fun (part, capacity) =>
      "f" ++ toString capacity ++ ":" ++ String.intercalate "," (part.map toString))
