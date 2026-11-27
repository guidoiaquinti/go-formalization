import GoLean.Bytes.TrimLeftByte
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let cut :: spare :: bytes := fields |
    throw (IO.userError "expected cut byte, spare capacity, and bytes")
  let cut ← wordInput 8 cut
  let raw ← bytes.mapM (wordInput 8)
  match GoLean.Bytes.trimLeftByte raw cut spare.toNat with
  | none => return "nil"
  | some (trimmed, cap) =>
    return toString cap ++ ":" ++ String.intercalate "," (trimmed.map (fun byte => toString byte.toNat))
