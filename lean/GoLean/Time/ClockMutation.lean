import GoLean.Time.Location

namespace GoLean.Time

private def packedFlag : BitVec 64 := 0x8000000000000000
private def nanoMask : BitVec 64 := 0x3fffffff
private def packedSpan : Int := 8589934591
private def firstPacked : Int := 59453308800

/-- Go 1.27.0 time.go:237–246. Existing wall bits are retained by OR. -/
def clockSetMono (s : LocatedTime) (reading : BitVec 64) : LocatedTime :=
  if s.words.wall.getLsbD 63 then { s with words := ⟨s.words.wall, reading⟩ }
  else
    let sec := s.words.ext.toInt
    if sec < firstPacked || firstPacked + packedSpan < sec then s
    else
      let wall := s.words.wall ||| packedFlag ||| (BitVec.ofInt 64 (sec - firstPacked) <<< 30)
      { s with words := ⟨wall, reading⟩ }

private def addExtended (s : LocatedTime) (delta : BitVec 64) : LocatedTime :=
  let previous := s.words.ext
  let sum := previous + delta
  let positive := delta.toInt > 0
  let ext :=
    if (sum.toInt > previous.toInt) == positive then sum
    else if positive then 0x7fffffffffffffff
    else 0x8000000000000001
  { s with words := ⟨s.words.wall, ext⟩ }

/-- Go 1.27.0 time.go:192–214. Packed failure strips monotonic state before adding
    to the full signed wall seconds. The negative saturation endpoint is -MaxInt64. -/
def clockAddSec (s : LocatedTime) (delta : BitVec 64) : LocatedTime :=
  if s.words.wall.getLsbD 63 then
    let packedSeconds := (s.words.wall <<< 1) >>> 31
    let sum := packedSeconds + delta
    if 0 ≤ sum.toInt && sum.toInt ≤ packedSpan then
      let wall := (s.words.wall &&& nanoMask) ||| (sum <<< 30) ||| packedFlag
      { s with words := ⟨wall, s.words.ext⟩ }
    else addExtended (clockStripMono s) delta
  else addExtended s delta

theorem clockSetMono_existing (s : LocatedTime) (reading : BitVec 64)
    (marked : s.words.wall.getLsbD 63 = true) :
    (clockSetMono s reading).words = ⟨s.words.wall, reading⟩ := by
  simp [clockSetMono, marked]

theorem clockSetMono_rejectEarly (s : LocatedTime) (reading : BitVec 64)
    (unmarked : s.words.wall.getLsbD 63 = false) (early : s.words.ext.toInt < 59453308800) :
    clockSetMono s reading = s := by
  unfold clockSetMono
  rw [unmarked]
  simp [firstPacked, early]

theorem clockAddSec_negativeOverflow :
    (clockAddSec ⟨⟨17, 0x8000000000000000⟩, 0⟩ 0xffffffffffffffff).words.ext =
      0x8000000000000001 := by decide

theorem clockAddSec_exactMin :
    (clockAddSec ⟨⟨17, 0x8000000000000001⟩, 0⟩ 0xffffffffffffffff).words.ext =
      0x8000000000000000 := by decide

end GoLean.Time
