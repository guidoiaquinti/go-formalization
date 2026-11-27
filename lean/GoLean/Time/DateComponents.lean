namespace GoLean.Time

def dateToAbsoluteDays (year : BitVec 64) (month : Nat) (day : Int) : Nat :=
  let rollover : Nat := if month < 3 then 1 else 0
  let shiftedMonth := month + 12 * rollover
  let y := (year - BitVec.ofNat 64 rollover + 292277022400).toNat
  let monthDays := (979 * shiftedMonth - 2919) / 32
  let centuryDays := (BitVec.ofNat 64 (146097 * (y / 100))).toNat / 4
  (BitVec.ofInt 64 (Int.ofNat (centuryDays + 1461 * (y % 100) / 4 + monthDays) + day - 1)).toNat

def gregorianLeap (year : Int) : Bool :=
  year % 4 == 0 && (year % 100 != 0 || year % 400 == 0)

def daysBeforeMonth (month : Nat) : Int :=
  (214 * Int.ofNat month - 211).tdiv 7 - (if month ≥ 3 then 2 else 0)

def daysInMonth (month : Nat) (year : Int) : Nat :=
  if month == 2 then (if gregorianLeap year then 29 else 28)
  else 30 + (month + month / 8) % 2

def normalizeFields (wordSize : Nat) (hi lo base : Int) : Int × Int :=
  ((BitVec.ofInt wordSize (hi + lo / base)).toInt, lo % base)

theorem dateToAbsoluteDays_epoch : dateToAbsoluteDays (BitVec.ofInt 64 (-292277022400)) 3 1 = 0 :=
  by decide

theorem gregorianLeap_centuries : gregorianLeap (-100) = false ∧ gregorianLeap (-400) = true :=
  by decide

theorem daysBeforeMonth_march : daysBeforeMonth 3 = 59 := rfl

theorem daysInMonth_february : daysInMonth 2 1900 = 28 ∧ daysInMonth 2 2000 = 29 := by decide

theorem normalizeFields_borrow : normalizeFields 64 10 (-1) 60 = (9, 59) := by decide

end GoLean.Time
