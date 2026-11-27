import GoLean.Strings.TrimRightByte
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let cut :: bytes := fields | throw (IO.userError "expected the cut byte and input")
  let raw ← bytes.mapM (wordInput 8)
  let trimmed := GoLean.Strings.trimRightByte raw (← wordInput 8 cut)
  return "bytes:" ++ String.intercalate "," (trimmed.map (fun byte => toString byte.toNat))
