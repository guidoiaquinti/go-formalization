import GoLean.Strings.Clone
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let 0 :: input := fields | throw (IO.userError "expected zero marker and string bytes")
  let bytes ← input.mapM (wordInput 8)
  let (out, fresh) := GoLean.Strings.clone bytes
  return String.intercalate " " (toString fresh :: out.map (fun byte => toString byte.toNat))
