import Init.Data.BitVec.Lemmas

namespace GoLean.Archive.Zip

structure CivilClock where
  year : Int
  month : Int
  day : Int
  hour : Int
  minute : Int
  second : Int

private def leapYearsBefore (year : Int) : Int :=
  let prior := year - 1
  prior / 4 - prior / 100 + prior / 400

private def leap (year : Int) : Bool :=
  year % 4 == 0 && (year % 100 != 0 || year % 400 == 0)

private def monthOffsets : Array Int :=
  #[0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334]

private def unixSeconds (civil : CivilClock) : Int :=
  let month := civil.month - 1
  let y := civil.year + month / 12
  let m := month % 12
  let days := 365 * (y - 1970) + leapYearsBefore y - leapYearsBefore 1970 +
    monthOffsets[m.toNat]! + (if m ≥ 2 && leap y then 1 else 0) + civil.day - 1
  days * 86400 + civil.hour * 3600 + civil.minute * 60 + civil.second

/-- Decode packed fields and normalize the resulting UTC civil time. -/
def msDosTimeToTime (date clock : BitVec 16) : Int :=
  let d : Int := Int.ofNat date.toNat
  let t : Int := Int.ofNat clock.toNat
  unixSeconds ⟨1980 + d / 512, d / 32 % 16, d % 32, t / 2048, t / 32 % 64, t % 32 * 2⟩

/-- The caller supplies the local civil fields read from its time.Time. -/
def timeToMsDosTime (civil : CivilClock) : BitVec 16 × BitVec 16 :=
  (BitVec.ofInt 16 (civil.day + 32 * civil.month + 512 * (civil.year - 1980)),
   BitVec.ofInt 16 (civil.second / 2 + 32 * civil.minute + 2048 * civil.hour))

def fileHeaderModTime := msDosTimeToTime

structure HeaderTimes where
  seconds : Int
  nanoseconds : Nat
  dateWord : BitVec 16
  timeWord : BitVec 16

/-- UTC civil fields are the projection of the supplied instant after UTC(). -/
def fileHeaderSetModTime (seconds : Int) (nanoseconds : Nat) (utc : CivilClock) : HeaderTimes :=
  let (d, t) := timeToMsDosTime utc
  ⟨seconds, nanoseconds, d, t⟩

theorem msDosTime_epoch : msDosTimeToTime 33 0 = 315532800 := by decide

theorem timeToMsDosTime_odd :
    timeToMsDosTime ⟨1980, 1, 1, 0, 0, 1⟩ = (33, 0) := by decide

theorem fileHeaderModTime_epoch : fileHeaderModTime 33 0 = 315532800 := by decide

theorem fileHeaderSetModTime_keepsInstant (seconds : Int) (nanos : Nat) (clock : CivilClock) :
    (fileHeaderSetModTime seconds nanos clock).seconds = seconds ∧
    (fileHeaderSetModTime seconds nanos clock).nanoseconds = nanos := ⟨rfl, rfl⟩

end GoLean.Archive.Zip
