import GoLean.Strings.TrimSpace
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let _ :: bytes := fields | throw (IO.userError "expected the input marker and bytes")
  let raw ← bytes.mapM (wordInput 8)
  return "bytes:" ++ String.intercalate "," ((GoLean.Strings.trimSpace raw).map (fun byte => toString byte.toNat))
