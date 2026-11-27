import GoLean.Time.ClockFields
import GoDiff.Math.Bits.Tuple

open GoLean.Time

def main : IO Unit := runTuple fun row => do
  let [seconds, days, century, year, day] := row | throw (IO.userError "expected five time fields")
  let seconds ← wordInput 64 seconds
  let days ← wordInput 64 days
  let (hour, minute, second) := absoluteClock seconds
  let (month, monthDay) := absoluteMonthDay (BitVec.ofInt 32 day)
  let leap := if absoluteLeap century.toNat year.toNat then 1 else 0
  return s!"{absoluteDays seconds}|{hour},{minute},{second}|{absoluteWeekday days}|{absoluteJanFeb day.toNat}|{leap}|{month},{monthDay}"
