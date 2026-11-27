import GoLean.Time.TimeCompare

namespace GoLean.Time

private def signedWord (n : Int) : Int := (BitVec.ofInt 64 n).toInt

def clockUnixMilli (t : ClockWords) : Int :=
  signedWord (clockUnix t * 1000 + Int.ofNat (clockNanosecond t / 1000000))

def clockUnixMicro (t : ClockWords) : Int :=
  signedWord (clockUnix t * 1000000 + Int.ofNat (clockNanosecond t / 1000))

def clockUnixNano (t : ClockWords) : Int :=
  signedWord (clockUnix t * 1000000000 + Int.ofNat (clockNanosecond t))

/-- The constructor stores the signed int32 nanoseconds through a uint64 conversion. -/
def unixTime (seconds : BitVec 64) (nanos : BitVec 32) : ClockWords :=
  ⟨BitVec.ofInt 64 nanos.toInt, seconds + 62135596800⟩

def unixFromParts (seconds nanos : Int) : ClockWords :=
  unixTime (BitVec.ofInt 64 (seconds + nanos / 1000000000))
    (BitVec.ofInt 32 (nanos % 1000000000))

def unixFromMilli (milli : Int) : ClockWords :=
  unixTime (BitVec.ofInt 64 (milli / 1000)) (BitVec.ofInt 32 (milli % 1000 * 1000000))

def unixFromMicro (micro : Int) : ClockWords :=
  unixTime (BitVec.ofInt 64 (micro / 1000000)) (BitVec.ofInt 32 (micro % 1000000 * 1000))

private def preEpoch : ClockWords := ⟨999999999, 62135596799⟩

theorem clockUnixMilli_negative : clockUnixMilli preEpoch = -1 := by decide

theorem clockUnixMicro_negative : clockUnixMicro preEpoch = -1 := by decide

theorem clockUnixNano_negative : clockUnixNano preEpoch = -1 := by decide

theorem unixTime_localWords (s : BitVec 64) (n : BitVec 32) :
    (unixTime s n).ext = s + 62135596800 := rfl

theorem unixFromParts_negative :
    (unixFromParts 0 (-1)).wall = 999999999 ∧ (unixFromParts 0 (-1)).ext = 62135596799 := by decide

theorem unixFromMilli_negative :
    (unixFromMilli (-1)).wall = 999000000 ∧ (unixFromMilli (-1)).ext = 62135596799 := by decide

theorem unixFromMicro_negative :
    (unixFromMicro (-1)).wall = 999999000 ∧ (unixFromMicro (-1)).ext = 62135596799 := by decide

end GoLean.Time
