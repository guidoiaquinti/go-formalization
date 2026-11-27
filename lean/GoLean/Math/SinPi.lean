import GoLean.Math.Cos
import GoLean.Math.FloorCeil

namespace GoLean.Math

private def sinPiAngle (arm64 fused : Bool) (octant x : Nat) : Option Nat :=
  let scaled := fun difference =>
    BitVec.ofNat 64 (multiplyWords 0x400921fb54442d18 difference)
  if octant = 0 then sinWord arm64 fused (scaled x)
  else if octant ≤ 2 then cosWord arm64 fused (scaled (subtractWords 0x3fe0000000000000 x))
  else if octant ≤ 4 then sinWord arm64 fused (scaled (subtractWords 0x3ff0000000000000 x))
  else if octant ≤ 6 then
    (cosWord arm64 fused (scaled (subtractWords x 0x3ff8000000000000))).map negateWord
  else sinWord arm64 fused (scaled (subtractWords x 0x4000000000000000))

private def sinPiReduce (word : BitVec 64) : Option (Nat × Nat) :=
  let x := word.toNat
  if floorWord word ≠ x then
    (modWords word (BitVec.ofNat 64 0x4000000000000000)).map fun reduced =>
      (reduced, 4 * magnitudeUnits reduced / 2 ^ 1074)
  else if 0x4340000000000000 ≤ x then some (0, 0)
  else
    let encoded := if x < 0x4330000000000000 then addWords x 0x4330000000000000 else x
    if encoded % 2 = 0 then some (0, 0) else some (0x3ff0000000000000, 4)

/-- Go 1.27.0 src/math/lgamma.go:325–369. The helper receives the
    nonnegative magnitude of the original negative Lgamma input.
    Accept both zero words and nonnegative finite inputs. -/
def sinPiWord (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  let magnitude := x % 2 ^ 63
  if magnitude = 0 then some (negateWord x)
  else if 0x7ff0000000000000 ≤ x then none
  else if x < 0x3fd0000000000000 then
    (sinWord arm64 fused (BitVec.ofNat 64 (multiplyWords 0x400921fb54442d18 x))).map negateWord
  else do
    let (reduced, octant) ← sinPiReduce word
    return negateWord (← sinPiAngle arm64 fused octant reduced)

theorem sinPiWord_zeros (arm64 fused : Bool) (word : BitVec 64)
    (zero : word.toNat % 2 ^ 63 = 0) :
    sinPiWord arm64 fused word = some (negateWord word.toNat) := by
  simp [sinPiWord, zero]

theorem sinPiWord_excluded (arm64 fused : Bool) (word : BitVec 64)
    (nonzero : word.toNat % 2 ^ 63 ≠ 0) (outside : 0x7ff0000000000000 ≤ word.toNat) :
    sinPiWord arm64 fused word = none := by
  simp [sinPiWord, nonzero, outside]

set_option exponentiation.threshold 4096 in
set_option maxRecDepth 4096 in
theorem sinPiWord_integer_boundaries (arm64 fused : Bool) :
    sinPiWord arm64 fused (BitVec.ofNat 64 0x4330000000000000) = some 0x8000000000000000 ∧
      sinPiWord arm64 fused (BitVec.ofNat 64 0x4330000000000001) = some 0x8000000000000000 ∧
      sinPiWord arm64 fused (BitVec.ofNat 64 0x4340000000000000) = some 0x8000000000000000 := by
  cases arm64 <;> cases fused <;> decide

end GoLean.Math
