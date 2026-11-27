import GoLean.Math.FloorCeil
import GoLean.Math.IsNegInt
import GoLean.Math.Sin
import GoLean.Math.Stirling

namespace GoLean.Math

private def gammaOne : Nat := 0x3ff0000000000000
private def gammaTiny : Nat := 0x3e112e0be826d695

private def gammaDown : Nat → Nat → Nat → Nat × Nat
  | 0, x, z => (x, z)
  | fuel + 1, x, z =>
    if 0x4008000000000000 ≤ x ∧ x < 2 ^ 63 then
      let x := subtractWords x gammaOne
      gammaDown fuel x (multiplyWords z x)
    else (x, z)

private def gammaNegativeUp : Nat → Nat → Nat → Nat × Nat
  | 0, x, z => (x, z)
  | fuel + 1, x, z =>
    if 2 ^ 63 + gammaTiny ≤ x then
      gammaNegativeUp fuel (addWords x gammaOne) (divideWords z x)
    else (x, z)

private def gammaPositiveUp : Nat → Nat → Nat → Nat × Nat
  | 0, x, z => (x, z)
  | fuel + 1, x, z =>
    if gammaTiny ≤ x ∧ x < 0x4000000000000000 then
      gammaPositiveUp fuel (addWords x gammaOne) (divideWords z x)
    else (x, z)

private def gammaPolynomial (fused : Bool) (x start : Nat) (rest : List Nat) : Nat :=
  rest.foldl (fun acc coefficient => multiplyAddWords fused x acc coefficient) start

private def gammaCentral (arm64 compileFused : Bool) (word : Nat) : Nat :=
  -- This path has |x| ≤ 33, so 34 recurrence steps suffice before reaching the central interval or the tiny cutoff.
  let (x, z) := gammaDown 34 word gammaOne
  let (x, z) := gammaNegativeUp 34 x z
  let (x, z) := gammaPositiveUp 2 x z
  if x = 0x4000000000000000 then z
  else if x % 2 ^ 63 < gammaTiny then
    if x % 2 ^ 63 = 0 then 0x7ff0000000000000
    else
      let correction := multiplyAddWords (arm64 || compileFused) 0x3fe2788cfc6fb619 x gammaOne
      divideWords z (multiplyWords correction x)
  else
    let x := subtractWords x 0x4000000000000000
    let p := gammaPolynomial arm64 x 0x3f24fcb839982153 [
      0x3f5384e3e686bfab, 0x3f8557cde9db14b0, 0x3fa863d918c423d3, 0x3fca8da9dcae7d31,
      0x3fdfa1373993e312, 0x3ff0000000000000]
    let q := gammaPolynomial arm64 x 0xbef8487a8400d3af [
      0x3f41ae8a29152573, 0xbf7240e4e750b44a, 0x3f8831ed5b1bb117, 0x3fa25779e33fde67,
      0xbfce071a9d4287c2, 0x3fb24944c9cd3c51, 0x3ff0000000000000]
    divideWords (multiplyWords z p) q

private def gammaReflection (arm64 compileFused expFused : Bool) (q : BitVec 64) : Option Nat := do
  let magnitude := q.toNat
  let whole := floorWord q
  let negative := (magnitudeUnits whole / 2 ^ 1074) % 2 = 0
  let signed := fun result => if negative then negateWord result else result
  let fraction := subtractWords magnitude whole
  let fraction := if 0x3fe0000000000000 < fraction then
      subtractWords magnitude (addWords whole gammaOne)
    else fraction
  let sine ← sinWord arm64 compileFused (BitVec.ofNat 64 (multiplyWords 0x400921fb54442d18 fraction))
  let z := multiplyWords magnitude sine % 2 ^ 63
  if z = 0 then return signed 0x7ff0000000000000
  let (first, second) ← stirlingWords arm64 compileFused expFused q
  if first = 0x7ff0000000000000 then return signed 0
  let firstProduct := multiplyWords z first
  let denominator := if firstProduct = 0x7ff0000000000000 then firstProduct else multiplyWords firstProduct second
  let result := if denominator = 0x7ff0000000000000 then
      divideWords (divideWords (divideWords 0x400921fb54442d18 z) first) second
    else divideWords 0x400921fb54442d18 denominator
  return signed result

/-- Go 1.27.0 src/math/gamma.go:131–214. The central recurrence reduces
    the argument without replacing the polynomial operations. The
    reflected domain retains two Stirling factors and the divide order
    used when their product overflows. None denotes NaN by class. -/
def gammaWord (arm64 compileFused expFused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  let magnitude := x % 2 ^ 63
  if isNegInt word ∨ x = 0xfff0000000000000 ∨ 0x7ff0000000000000 < magnitude then none
  else if x = 0x7ff0000000000000 then some x
  else if magnitude = 0 then some (if 2 ^ 63 ≤ x then 0xfff0000000000000 else 0x7ff0000000000000)
  else if 0x4040800000000000 < magnitude then
    if 2 ^ 63 ≤ x then gammaReflection arm64 compileFused expFused (BitVec.ofNat 64 magnitude)
    else do
      let (first, second) ← stirlingWords arm64 compileFused expFused word
      return if first = 0x7ff0000000000000 then first else multiplyWords first second
  else some (gammaCentral arm64 compileFused x)

theorem gammaWord_negativeInteger (arm64 compileFused expFused : Bool) (word : BitVec 64)
    (negativeInteger : isNegInt word = true) :
    gammaWord arm64 compileFused expFused word = none := by
  simp [gammaWord, negativeInteger]

theorem gammaWord_specials (arm64 compileFused expFused : Bool) :
    gammaWord arm64 compileFused expFused 0 = some 0x7ff0000000000000 ∧
      gammaWord arm64 compileFused expFused (BitVec.ofNat 64 0x8000000000000000) =
        some 0xfff0000000000000 ∧
      gammaWord arm64 compileFused expFused (BitVec.ofNat 64 0x7ff0000000000000) =
        some 0x7ff0000000000000 ∧
      gammaWord arm64 compileFused expFused (BitVec.ofNat 64 0xfff0000000000000) = none := by
  cases arm64 <;> cases compileFused <;> cases expFused <;> decide

end GoLean.Math
