import GoLean.Bytes.TrimLeft
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let cutsetLength :: spare :: bytes := fields |
    throw (IO.userError "expected cutset length, spare capacity, and bytes")
  let raw ← bytes.mapM (wordInput 8)
  let cutsetLength := cutsetLength.toNat
  match GoLean.Bytes.trimLeft (raw.drop cutsetLength) (raw.take cutsetLength) spare.toNat with
  | none => return "nil"
  | some (trimmed, capacity) =>
    return toString capacity ++ ":" ++ String.intercalate "," (trimmed.map (fun byte => toString byte.toNat))
