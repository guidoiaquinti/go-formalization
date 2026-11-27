import GoLean.Unicode.Utf16.AppendRune
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let r :: words := fields | throw (IO.userError "expected a rune and a sequence of words")
  let dst ← words.mapM fun word => do
    return (← wordInput 16 word).toNat
  return String.intercalate " " ((GoLean.Unicode.Utf16.appendRune dst r).map toString)
