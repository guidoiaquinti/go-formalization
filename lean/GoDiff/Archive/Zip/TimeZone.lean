import GoLean.Archive.Zip.TimeZone
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let [offset] := row | throw (IO.userError "expected signed nanoseconds")
  let result := GoLean.Archive.Zip.timeZone offset
  return s!"{result.seconds} {result.cached}"
