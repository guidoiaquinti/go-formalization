import GoLean.Unicode.Utf8.DecodeRuneInString
import GoDiff.Math.Bits.Tuple

open GoLean.Unicode.Utf8

def main : IO Unit := runTuple fun fields => do
  let 0 :: words := fields | throw (IO.userError "expected marker and bytes")
  let bytes ← words.mapM fun n => return (← wordInput 8 n).toNat
  let (rune, width) := decodeRuneInString bytes
  return s!"{rune} {width}"
