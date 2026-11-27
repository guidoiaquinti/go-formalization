import GoLean.Time.DateComponents
import GoDiff.Math.Bits.Tuple

open GoLean.Time

def main : IO Unit := runTuple fun row => do
  let [w, year, month, day, hi, lo, base] := row | throw (IO.userError "expected native date fields")
  let (h, l) := normalizeFields w.toNat hi lo base
  return s!"{dateToAbsoluteDays (BitVec.ofInt 64 year) month.toNat day}|{gregorianLeap year}|{daysBeforeMonth month.toNat}|{daysInMonth month.toNat year}|{h},{l}"
