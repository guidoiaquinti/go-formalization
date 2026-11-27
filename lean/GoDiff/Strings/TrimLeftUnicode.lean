import GoLean.Strings.Trim
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let cutsetLength :: bytes := fields | throw (IO.userError "expected cutset length and bytes")
  let cutsetLength := cutsetLength.toNat
  let raw ← bytes.mapM (wordInput 8)
  let trimmed := GoLean.Strings.trimLeftUnicode (raw.drop cutsetLength) (raw.take cutsetLength)
  return "bytes:" ++ String.intercalate "," (trimmed.map (fun byte => toString byte.toNat))
