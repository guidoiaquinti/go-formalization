import GoLean.Time.ClockFields

namespace GoLean.Time

def splitAbsoluteDays (days : BitVec 64) : Nat × Nat × Nat :=
  let scaled := (4 * days + 3).toNat
  let century := scaled / 146097
  let inCentury := (scaled % 146097 / 4) * 4 + 3
  let product := 2939745 * inCentury
  (century, product / 4294967296, product % 4294967296 / 2939745 / 4)

def civilYear (wordSize century year janFeb : Nat) : Int :=
  (BitVec.ofInt wordSize (Int.ofNat (century * 100 + year + janFeb) - 292277022400)).toInt

def civilMonth (month janFeb : Nat) : Int := Int.ofNat month - 12 * Int.ofNat janFeb

def civilYearDay (day janFeb : Nat) (leap : Bool) : Int :=
  Int.ofNat day + 60 + (if leap && janFeb == 0 then 1 else 0) - 365 * Int.ofNat janFeb

def dateFromAbsoluteDays (wordSize : Nat) (days : BitVec 64) : Int × Int × Nat :=
  let (century, year, yday) := splitAbsoluteDays days
  let (month, mday) := absoluteMonthDay (BitVec.ofNat 32 yday)
  let janFeb := absoluteJanFeb yday
  (civilYear wordSize century year janFeb, civilMonth month janFeb, mday)

def yearDayFromAbsoluteDays (wordSize : Nat) (days : BitVec 64) : Int × Int :=
  let (century, year, day) := splitAbsoluteDays days
  let janFeb := absoluteJanFeb day
  (civilYear wordSize century year janFeb, civilYearDay day janFeb (absoluteLeap century year))

theorem splitAbsoluteDays_cycle : splitAbsoluteDays 146097 = (4, 0, 0) := by decide

theorem civilYear_rollover : civilYear 64 2922770224 99 1 = 100 := by decide

theorem civilMonth_feb : civilMonth 14 1 = 2 := rfl

theorem civilYearDay_jan : civilYearDay 306 1 true = 1 := rfl

theorem dateFromAbsoluteDays_rollover : dateFromAbsoluteDays 64 306 = (-292277022399, 1, 1) := by decide

theorem yearDayFromAbsoluteDays_rollover : yearDayFromAbsoluteDays 64 306 = (-292277022399, 1) := by decide

end GoLean.Time
