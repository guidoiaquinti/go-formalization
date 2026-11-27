import GoLean.Time.EpochUnits
import GoDiff.Math.Bits.Tuple

open GoLean.Time

private def wire (s : ClockWords) := s!"{s.wall.toNat},{s.ext.toNat}"

def main : IO Unit := runTuple fun row => do
  let [wall, ext, seconds, nsec, units] := row | throw (IO.userError "expected raw time and epoch inputs")
  let value : ClockWords := ⟨← wordInput 64 wall, ← wordInput 64 ext⟩
  return String.intercalate "|" [toString (clockUnixMilli value), toString (clockUnixMicro value),
    toString (clockUnixNano value), wire (unixTime (BitVec.ofInt 64 seconds) (BitVec.ofInt 32 nsec)),
    wire (unixFromParts seconds nsec), wire (unixFromMilli units), wire (unixFromMicro units)]
