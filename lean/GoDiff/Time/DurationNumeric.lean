import GoLean.Time.DurationNumeric
import GoDiff.Math.Bits.Tuple
open GoLean.Time

def main : IO Unit := runTuple fun row => do
  let [d, unit] := row | throw (IO.userError "expected duration and divisor")
  return String.intercalate "|" [
    toString (durationNanoseconds d), toString (durationMicroseconds d), toString (durationMilliseconds d),
    toString (durationSeconds d).toBits, toString (durationMinutes d).toBits, toString (durationHours d).toBits,
    toString (durationTruncate d unit), toString (durationRound d unit), toString (durationAbs d)]
