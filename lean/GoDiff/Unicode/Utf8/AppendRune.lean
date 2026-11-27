import GoLean.Unicode.Utf8.AppendRune
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let r :: words := fields | throw (IO.userError "expected rune and destination bytes")
  let dst ← words.mapM fun n => return (← wordInput 8 n).toNat
  return String.intercalate " " ((GoLean.Unicode.Utf8.appendRune dst r).map toString)
