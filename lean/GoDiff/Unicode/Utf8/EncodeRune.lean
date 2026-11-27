import GoLean.Unicode.Utf8.EncodeRune
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let r :: words := fields | throw (IO.userError "expected rune and destination bytes")
  let dst ← words.mapM fun n => return (← wordInput 8 n).toNat
  let (written, result) := GoLean.Unicode.Utf8.encodeRune dst r
  let status := written.map toString |>.getD "bounds"
  return String.intercalate " " (status :: result.map toString)
