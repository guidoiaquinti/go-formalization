import GoLean.Unicode.Utf8.RuneCount
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let 0 :: input := fields | throw (IO.userError "expected zero marker and input bytes")
  let words ← input.mapM fun n => do
    return (← wordInput 8 n).toNat
  return s!"{GoLean.Unicode.Utf8.runeCount words} {GoLean.Unicode.Utf8.runeCountInString words}"
