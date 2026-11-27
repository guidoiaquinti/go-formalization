import Init.Data.BitVec.Lemmas

namespace GoLean.Time

structure ClockWords where
  wall : BitVec 64
  ext : BitVec 64

def clockNsec (t : ClockWords) : Nat := (t.wall &&& (2^30 - 1)).toNat

def clockSec (t : ClockWords) : Int :=
  if t.wall.getLsbD 63 then
    59453308800 + ((t.wall <<< 1) >>> 31).toNat
  else t.ext.toInt

def clockUnixSec (t : ClockWords) : Int :=
  (BitVec.ofInt 64 (clockSec t - 62135596800)).toInt

def clockMono (t : ClockWords) : Int :=
  if t.wall.getLsbD 63 then t.ext.toInt else 0

def clockNanosecond (t : ClockWords) := clockNsec t

def clockUnix (t : ClockWords) := clockUnixSec t

def clockIsZero (t : ClockWords) : Bool := t.wall == 0 && t.ext == 0

private def compareKeys (t u : ClockWords) : Int × Int :=
  if t.wall.getLsbD 63 && u.wall.getLsbD 63 then (t.ext.toInt, u.ext.toInt)
  else
    let ts := clockSec t
    let us := clockSec u
    if ts == us then (clockNsec t, clockNsec u) else (ts, us)

def clockBefore (t u : ClockWords) : Bool :=
  let (a,b) := compareKeys t u
  decide (a < b)

def clockAfter (t u : ClockWords) : Bool :=
  let (a,b) := compareKeys t u
  decide (a > b)

def clockEqual (t u : ClockWords) : Bool :=
  let (a,b) := compareKeys t u
  a == b

def clockCompare (t u : ClockWords) : Int :=
  let (a,b) := compareKeys t u
  if a < b then -1 else if a > b then 1 else 0

theorem clockNsec_maxPacked : clockNsec ⟨0xffffffffffffffff, 0⟩ = 1073741823 := rfl

theorem clockSec_withoutMono (t : ClockWords) (h : t.wall.getLsbD 63 = false) :
    clockSec t = t.ext.toInt := by rw [clockSec, h]; rfl

theorem clockUnixSec_epoch : clockUnixSec ⟨0, 62135596800⟩ = 0 := rfl

theorem clockMono_withoutMono (t : ClockWords) (h : t.wall.getLsbD 63 = false) :
    clockMono t = 0 := by rw [clockMono, h]; rfl

theorem clockNanosecond_nsec (t : ClockWords) : clockNanosecond t = clockNsec t := rfl

theorem clockUnix_seconds (t : ClockWords) : clockUnix t = clockUnixSec t := rfl

theorem clockIsZero_words (t : ClockWords) :
    clockIsZero t = true ↔ t.wall = 0 ∧ t.ext = 0 := by simp [clockIsZero]

theorem clockBefore_mono (t u : ClockWords) (ht : t.wall.getLsbD 63 = true)
    (hu : u.wall.getLsbD 63 = true) :
    clockBefore t u = decide (t.ext.toInt < u.ext.toInt) := by simp [clockBefore, compareKeys, ht, hu]

theorem clockAfter_mono (t u : ClockWords) (ht : t.wall.getLsbD 63 = true)
    (hu : u.wall.getLsbD 63 = true) :
    clockAfter t u = decide (t.ext.toInt > u.ext.toInt) := by simp [clockAfter, compareKeys, ht, hu]

theorem clockEqual_mono (t u : ClockWords) (ht : t.wall.getLsbD 63 = true)
    (hu : u.wall.getLsbD 63 = true) :
    clockEqual t u = (t.ext.toInt == u.ext.toInt) := by simp [clockEqual, compareKeys, ht, hu]

theorem clockCompare_mono (t u : ClockWords) (ht : t.wall.getLsbD 63 = true)
    (hu : u.wall.getLsbD 63 = true) :
    clockCompare t u = (if t.ext.toInt < u.ext.toInt then -1 else if t.ext.toInt > u.ext.toInt then 1 else 0) := by
  simp [clockCompare, compareKeys, ht, hu]

end GoLean.Time
