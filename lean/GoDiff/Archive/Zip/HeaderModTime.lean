import GoLean.Archive.Zip.HeaderModTime
import GoDiff.Math.Bits.Tuple

open GoLean.Archive.Zip

def main : IO Unit := runTuple fun row => do
  let [wall, ext, d, t] := row | throw (IO.userError "expected clock and DOS words")
  let wall ← wordInput 64 wall
  let ext ← wordInput 64 ext
  let date ← wordInput 16 d
  let clock ← wordInput 16 t
  let result := headerFileInfoModTime ⟨wall, ext⟩ date clock
  return s!"{result.seconds}:{result.nanoseconds}"
