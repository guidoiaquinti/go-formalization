import Init.Data.OfScientific
import Init.Data.Int.DivMod

namespace GoLean.Time

def durationNanoseconds (d : Int) : Int := d
def durationMicroseconds (d : Int) : Int := d.tdiv 1000
def durationMilliseconds (d : Int) : Int := d.tdiv 1000000

/-- Go 1.27.0 src/time/time.go:1090–1109. Round each floating operation
at the source-prescribed scale after splitting the signed duration. -/
private def fraction (nanos unit : Int) : Float :=
  Float.ofInt (nanos.tdiv unit) + Float.ofInt (nanos.tmod unit) / Float.ofInt unit

def durationSeconds (d : Int) : Float := fraction d 1000000000
def durationMinutes (d : Int) : Float := fraction d 60000000000
def durationHours (d : Int) : Float := fraction d 3600000000000

def durationTruncate (d m : Int) : Int :=
  if m ≤ 0 then d else d - d.tmod m

/-- Go 1.27.0 src/time/time.go:1131–1153. Overflow detection on the
signed source implements saturation of the nearest mathematical multiple. -/
def durationRound (d m : Int) : Int :=
  if m ≤ 0 then d else
    let r := d.tmod m
    let toward := d - r
    if Int.ofNat r.natAbs * 2 < m then toward
    else if d < 0 then max (-9223372036854775808) (toward - m)
    else min 9223372036854775807 (toward + m)

def durationAbs (d : Int) : Int := min 9223372036854775807 (Int.ofNat d.natAbs)

theorem durationNanoseconds_identity (d : Int) : durationNanoseconds d = d := rfl
theorem durationMicroseconds_negativeBoundary : durationMicroseconds (-999) = 0 := rfl
theorem durationMilliseconds_negativeBoundary : durationMilliseconds (-999999) = 0 := rfl
theorem durationSeconds_split (d : Int) :
    durationSeconds d =
      Float.ofInt (d.tdiv 1000000000) + Float.ofInt (d.tmod 1000000000) / Float.ofInt 1000000000 := rfl
theorem durationMinutes_split (d : Int) :
    durationMinutes d =
      Float.ofInt (d.tdiv 60000000000) + Float.ofInt (d.tmod 60000000000) / Float.ofInt 60000000000 := rfl
theorem durationHours_split (d : Int) :
    durationHours d =
      Float.ofInt (d.tdiv 3600000000000) + Float.ofInt (d.tmod 3600000000000) / Float.ofInt 3600000000000 := rfl
theorem durationTruncate_nonpositive (d m : Int) (h : m ≤ 0) : durationTruncate d m = d := by
  simp [durationTruncate, h]
theorem durationRound_nonpositive (d m : Int) (h : m ≤ 0) : durationRound d m = d := by
  simp [durationRound, h]
theorem durationRound_ties : durationRound (-3) 2 = -4 ∧ durationRound 3 2 = 4 := ⟨rfl, rfl⟩
theorem durationRound_saturates :
    durationRound 9223372036854775807 2 = 9223372036854775807 ∧
      durationRound (-9223372036854775807) 4 = -9223372036854775808 := ⟨rfl, rfl⟩

theorem durationAbs_min : durationAbs (-9223372036854775808) = 9223372036854775807 := rfl

end GoLean.Time
