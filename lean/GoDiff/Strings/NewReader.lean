import GoLean.Strings.NewReader
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let 0 :: bytes := fields | throw (IO.userError "expected marker and string bytes")
  let s := GoLean.Strings.newReader (← bytes.mapM (wordInput 8))
  return String.intercalate " " (toString s.cursor :: toString s.previousRune ::
    s.contents.map (fun byte => toString byte.toNat))
