import GoLean.Strings.Title
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let _ :: raw := row | throw (IO.userError "expected marker and bytes")
  let input ← raw.mapM (wordInput 8)
  let output := GoLean.Strings.title input
  return "bytes:" ++ String.intercalate "," (output.map toString)
