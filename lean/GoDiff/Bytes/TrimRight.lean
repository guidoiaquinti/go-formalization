import GoLean.Bytes.TrimRight
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let cutsetLength :: spare :: isNil :: bytes := fields |
    throw (IO.userError "expected cutset length, spare capacity, nil marker, and bytes")
  let raw ← bytes.mapM (wordInput 8)
  let cutsetLength := cutsetLength.toNat
  let input := raw.drop cutsetLength
  let slice := if isNil == 1 then none else some (input, input.length + spare.toNat)
  match GoLean.Bytes.trimRight slice (raw.take cutsetLength) with
  | none => return "nil"
  | some (trimmed, capacity) =>
    return toString capacity ++ ":" ++ String.intercalate "," (trimmed.map (fun byte => toString byte.toNat))
