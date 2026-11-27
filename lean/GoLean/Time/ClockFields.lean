namespace GoLean.Time

def absoluteDays (seconds : BitVec 64) : Nat := seconds.toNat / 86400

def absoluteClock (seconds : BitVec 64) : Nat × Nat × Nat :=
  let withinDay := seconds.toNat % 86400
  (withinDay / 3600, withinDay % 3600 / 60, withinDay % 60)

def absoluteWeekday (days : BitVec 64) : Nat := (days + 3).toNat % 7

def absoluteJanFeb (day : Nat) : Nat := if day ≥ 306 then 1 else 0

def absoluteLeap (century year : Nat) : Bool :=
  year % 4 == 0 && (year != 0 || century % 4 == 0)

def absoluteMonthDay (day : BitVec 32) : Nat × Nat :=
  let word := (2141 * day + 197913).toNat
  (word / 65536, 1 + word % 65536 / 2141)

theorem absoluteDays_firstBoundary : absoluteDays 86399 = 0 ∧ absoluteDays 86400 = 1 := by decide

theorem absoluteClock_lastSecond : absoluteClock 86399 = (23, 59, 59) := by decide

theorem absoluteWeekday_wrap : absoluteWeekday 18446744073709551615 = 2 := by decide

theorem absoluteJanFeb_boundary : absoluteJanFeb 305 = 0 ∧ absoluteJanFeb 306 = 1 := by decide

theorem absoluteLeap_centuries : absoluteLeap 1 0 = false ∧ absoluteLeap 4 0 = true := by decide

theorem absoluteMonthDay_monthBoundary :
    absoluteMonthDay 30 = (3, 31) ∧ absoluteMonthDay 31 = (4, 1) := by decide

end GoLean.Time
