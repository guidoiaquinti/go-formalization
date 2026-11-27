import GoLean.Time.CalendarConversion
import GoDiff.Math.Bits.Tuple

open GoLean.Time

def main : IO Unit := runTuple fun row => do
  let [w, d, c, y, j, m, yd, l] := row | throw (IO.userError "expected native width and calendar fields")
  let d ← wordInput 64 d
  let w := w.toNat
  let j := j.toNat
  let (century, year, day) := splitAbsoluteDays d
  let (dateYear, month, dateDay) := dateFromAbsoluteDays w d
  let (dayYear, yearDay) := yearDayFromAbsoluteDays w d
  return s!"{century},{year},{day}|{civilYear w c.toNat y.toNat j}|{civilMonth m.toNat j}|{civilYearDay yd.toNat j (l != 0)}|{dateYear},{month},{dateDay}|{dayYear},{yearDay}"
