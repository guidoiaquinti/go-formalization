import GoLean.Strings.Lines
import GoDiff.Math.Bits.Tuple

private def encode (batches : List (List (List (BitVec 8)))) : String :=
  let describeLine := fun (line : List (BitVec 8)) =>
    String.intercalate "," (line.map (fun byte => toString byte.toNat))
  String.intercalate "|" (batches.map (fun batch => String.intercalate ";" (batch.map describeLine)))

def main : IO Unit := runTuple fun fields => do
  let first :: second :: input := fields | throw (IO.userError "expected two budgets and string bytes")
  let bytes ← input.mapM (wordInput 8)
  let (firstBatch, rest) := GoLean.Strings.lines bytes first.toNat
  let (secondBatch, rest) := GoLean.Strings.lines rest second.toNat
  let (lastBatch, rest) := GoLean.Strings.lines rest (bytes.length + 1)
  let (after, _) := GoLean.Strings.lines rest 1
  return encode [firstBatch, secondBatch, lastBatch, after]
