import GoLean.Archive.Zip.DetectUTF8
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let 0 :: words := fields | throw (IO.userError "expected zero tag and bytes")
  let bytes ← words.mapM fun n => return (← wordInput 8 n).toNat
  let (valid, require) := GoLean.Archive.Zip.detectUTF8 bytes
  return s!"{valid} {require}"
