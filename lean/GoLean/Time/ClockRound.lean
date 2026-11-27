import GoLean.Time.ClockArithmetic

namespace GoLean.Time

private def billion : Int := 1000000000

private def i64 (value : Int) : Int := (BitVec.ofInt 64 value).toInt

def lessThanHalf (a b : BitVec 64) : Bool := (a + a).toNat < b.toNat

/-- The private helper has a positive-divisor contract. The general branch
    represents the long division by unbounded arithmetic on its unsigned input. -/
def clockDiv (words : ClockWords) (d : Int) : Int × Int :=
  let sec := clockSec words
  let nano := Int.ofNat (clockNsec words)
  let neg := sec < 0
  let (s, n) := if neg then
    if nano = 0 then (i64 (-sec), 0) else (i64 (-sec - 1), billion - nano)
    else (sec, nano)
  let (q, r) :=
    if d < billion && billion % (2 * d) = 0 then (n.tdiv d, n.tmod d)
    else if d % billion = 0 then
      let unit := d / billion
      (s.tdiv unit, s.tmod unit * billion + n)
    else
      let combined := Int.ofNat (BitVec.ofInt 64 s).toNat * billion +
        Int.ofNat (BitVec.ofInt 64 n).toNat
      (combined / d, combined % d)
  if neg && r != 0 then (1 - q % 2, i64 (d - r)) else (q % 2, i64 r)

def clockTruncate (s : LocatedTime) (duration : BitVec 64) : LocatedTime :=
  let s := clockStripMono s
  let d := duration.toInt
  if d ≤ 0 then s
  else clockAdd s (BitVec.ofInt 64 (- (clockDiv s.words d).2))

def clockRound (s : LocatedTime) (duration : BitVec 64) : LocatedTime :=
  let s := clockStripMono s
  let d := duration.toInt
  if d ≤ 0 then s
  else
    let r := (clockDiv s.words d).2
    let delta := if lessThanHalf (BitVec.ofInt 64 r) duration then -r else d - r
    clockAdd s (BitVec.ofInt 64 delta)

theorem lessThanHalf_equality : lessThanHalf 5 10 = false := rfl

theorem clockDiv_negative :
    clockDiv ⟨999999999, BitVec.ofInt 64 (-1)⟩ 1000000000 = (1, 999999999) := by decide

theorem clockTruncate_zero (s : LocatedTime) :
    clockTruncate s 0 = clockStripMono s := by simp [clockTruncate]

theorem clockRound_zero (s : LocatedTime) :
    clockRound s 0 = clockStripMono s := by simp [clockRound]

end GoLean.Time
