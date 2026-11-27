import GoLean.Bytes.SplitAfterSeq
import GoDiff.Math.Bits.Tuple

private def describe (parts : List (List Nat × Nat)) : String :=
  String.intercalate ";" (parts.map fun part =>
    "f" ++ toString part.2 ++ ":" ++ String.intercalate "," (part.1.map toString))

def main : IO Unit := runTuple fun row => do
  let separatorSize :: first :: second :: raw := row |
    throw (IO.userError "expected separator size, budgets and bytes")
  let bytes ← raw.mapM (wordInput 8)
  let values := bytes.map BitVec.toNat
  let separator := values.take separatorSize.toNat
  let input := values.drop separatorSize.toNat
  let (firstBatch, rest) := GoLean.Bytes.splitAfterSeq input separator first.toNat
  let (secondBatch, rest) := GoLean.Bytes.splitAfterSeq rest separator second.toNat
  let (full, rest) := GoLean.Bytes.splitAfterSeq rest separator (input.length + 2)
  let (after, _) := GoLean.Bytes.splitAfterSeq rest separator 1
  return String.intercalate "|" ([firstBatch, secondBatch, full, after].map describe)
