import GoLean.Bytes.TrimRightByte
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let cut :: capacity :: isNil :: bytes := fields |
    throw (IO.userError "expected cut byte, capacity, nil marker, and bytes")
  let cut ← wordInput 8 cut
  let raw ← bytes.mapM (wordInput 8)
  let input := if isNil == 1 then none else some (raw, capacity.toNat)
  match GoLean.Bytes.trimRightByte input cut with
  | none => return "nil"
  | some (trimmed, cap) =>
    return toString cap ++ ":" ++ String.intercalate "," (trimmed.map (fun byte => toString byte.toNat))
