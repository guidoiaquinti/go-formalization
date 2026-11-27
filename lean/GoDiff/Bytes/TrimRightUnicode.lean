import GoLean.Bytes.TrimRightUnicode
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let cutsetLength :: capacity :: bytes := fields |
    throw (IO.userError "expected cutset length, capacity or nil marker, and bytes")
  let raw ← bytes.mapM (wordInput 8)
  let cutsetLength := cutsetLength.toNat
  let input := if capacity < 0 then none else some (raw.drop cutsetLength, capacity.toNat)
  match GoLean.Bytes.trimRightUnicode input (raw.take cutsetLength) with
  | none => return "nil"
  | some (trimmed, cap) =>
    return toString cap ++ ":" ++ String.intercalate "," (trimmed.map (fun byte => toString byte.toNat))
