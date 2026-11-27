import GoLean.Time.ClockMutation

namespace GoLean.Time

private def billion : Int := 1000000000
private def nanosMask : BitVec 64 := 0x3fffffff
private def maxDuration : BitVec 64 := 0x7fffffffffffffff
private def minDuration : BitVec 64 := 0x8000000000000000

/-- Mutate wall and signed monotonic words in the order used by Go. -/
def clockAdd (s : LocatedTime) (duration : BitVec 64) : LocatedTime :=
  let d := duration.toInt
  let seconds := d.tdiv billion
  let nanos := Int.ofNat (clockNsec s.words) + d.tmod billion
  let (seconds, nanos) :=
    if nanos ≥ billion then (seconds + 1, nanos - billion)
    else if nanos < 0 then (seconds - 1, nanos + billion)
    else (seconds, nanos)
  let wall := (s.words.wall &&& ~~~nanosMask) ||| BitVec.ofInt 64 nanos
  let advanced := clockAddSec {s with words := ⟨wall, s.words.ext⟩} (BitVec.ofInt 64 seconds)
  if advanced.words.wall.getLsbD 63 then
    let old := advanced.words.ext
    let sum := old + duration
    if (d < 0 && sum.toInt > old.toInt) || (d > 0 && sum.toInt < old.toInt) then
      clockStripMono advanced
    else {advanced with words := ⟨advanced.words.wall, sum⟩}
  else advanced

def subMono (t u : BitVec 64) : BitVec 64 :=
  let delta := t - u
  if delta.toInt < 0 && t.toInt > u.toInt then maxDuration
  else if delta.toInt > 0 && t.toInt < u.toInt then minDuration
  else delta

def clockSub (t u : LocatedTime) : BitVec 64 :=
  if t.words.wall.getLsbD 63 && u.words.wall.getLsbD 63 then subMono t.words.ext u.words.ext
  else
    let seconds := clockSec t.words - clockSec u.words
    let nanos := Int.ofNat (clockNsec t.words) - Int.ofNat (clockNsec u.words)
    let duration := BitVec.ofInt 64 (seconds * billion + nanos)
    if clockEqual (clockAdd u duration).words t.words then duration
    else if clockBefore t.words u.words then minDuration
    else maxDuration

theorem clockAdd_borrow :
    (clockAdd ⟨⟨0, 2⟩, 5⟩ (-1)).words.ext = 1 ∧
      (clockAdd ⟨⟨0, 2⟩, 5⟩ (-1)).words.wall = 999999999 := ⟨rfl, rfl⟩

theorem clockSub_monotonic (t u : LocatedTime)
    (ht : t.words.wall.getLsbD 63 = true) (hu : u.words.wall.getLsbD 63 = true) :
    clockSub t u = subMono t.words.ext u.words.ext := by simp [clockSub, ht, hu]

theorem subMono_saturation :
    subMono 0x7fffffffffffffff 0x8000000000000000 = 0x7fffffffffffffff ∧
    subMono 0x8000000000000000 0x7fffffffffffffff = 0x8000000000000000 := ⟨rfl, rfl⟩

end GoLean.Time
