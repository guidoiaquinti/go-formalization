import GoLean.Bytes.TrimSpace
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let spare :: bytes := fields | throw (IO.userError "expected spare capacity and bytes")
  let raw ← bytes.mapM (wordInput 8)
  match GoLean.Bytes.trimSpace raw spare.toNat with
  | none => return "nil"
  | some (trimmed, capacity) =>
    return toString capacity ++ ":" ++ String.intercalate "," (trimmed.map (fun byte => toString byte.toNat))
