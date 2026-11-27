import GoLean.Math.Log
import GoLean.Math.SinPi

namespace GoLean.Math

private def lgammaOne : Nat := 0x3ff0000000000000
private def lgammaTwo : Nat := 0x4000000000000000
private def lgammaHalf : Nat := 0x3fe0000000000000

private def lgammaPolynomial (fused : Bool) (x first : Nat) (rest : List Nat) : Nat :=
  rest.foldl (fun acc coefficient => multiplyAddWords fused x acc coefficient) first

private def lgammaLowCorrection (arm64 fused : Bool) (which y : Nat) : Nat :=
  if which = 0 then
    let square := multiplyWords y y
    let even := lgammaPolynomial arm64 square 0x3efa7074428cfa52 [
    0x3f2cf2eced10e54d, 0x3f538a94116f3f5d, 0x3f7e404fb68fefe8, 0x3fb13e001a5562a7,
    0x3fb3c467e37db0c8]
    let odd := lgammaPolynomial arm64 square 0x3f07858e90a45837 [
    0x3f1c5088987dfb07, 0x3f40b6c689b99c00, 0x3f67add8ccb7926b, 0x3f951322ac92547b,
    0x3fd4a34cc4a60fad]
    let p :=
      if arm64 then multiplyAddWords true square odd (multiplyWords y even)
      else multiplyAddWords fused y even (multiplyWords square odd)
    multiplyAddWords arm64 y 0xbfe0000000000000 p
  else if which = 1 then
    let square := multiplyWords y y
    let cube := multiplyWords square y
    let p1 := lgammaPolynomial arm64 cube 0x3f34af6d6c0ebbf7 [
    0xbf56fe8ebf2d1af1, 0x3f78fce0e370e344, 0xbfa0c9a8df35b713, 0x3fdef72bc8ee38a2]
    let p2 := lgammaPolynomial arm64 cube 0xbf347f24ecc38c38 [
    0x3f4cdf0cef61a8e9, 0xbf6e2effb3e914d7, 0x3f9266e7970af9ec, 0xbfc2e4278dc6c509]
    let p3 := lgammaPolynomial arm64 cube 0x3f35fd3ee8c2d3f4 [
    0xbf41a6109c73e0ec, 0x3f6282d32e15c915, 0xbf851f9fba91ec6a, 0x3fb08b4294d5419b]
    let combined := multiplyAddWords fused y p3 p2
    let tail := multiplyAddWords arm64 (negateWord cube) combined 0xbc50c7caa48a971f
    let p := multiplyAddWords arm64 square p1 (negateWord tail)
    addWords 0xbfbf19b9bcc38a42 p
  else
    let p := multiplyWords y (lgammaPolynomial arm64 y 0x3f8b678bbf2bab09 [
    0x3fcd4eaef6010924, 0x3fef497644ea8450, 0x3ff7475cd119bd6f, 0x3fe4401e8b005dff,
    0xbfb3c467e37db0c8])
    let q := multiplyAddWords fused y (lgammaPolynomial arm64 y 0x3f6a5abb57d0cf61 [
    0x3fbaae55d6537c88, 0x3fe89dfbe45050af, 0x40010725a42b18f5, 0x4003a5d7c2bd619c]) lgammaOne
    multiplyAddWords fused y 0xbfe0000000000000 (divideWords p q)

private def lgammaLow (arm64 fused : Bool) (x : Nat) : Option Nat := do
  let (base, y, which) ←
    if x ≤ 0x3feccccccccccccd then
      (logWord arm64 (BitVec.ofNat 64 x)).map fun logarithm =>
        let base := negateWord logarithm
        if 0x3fe76987d0eaed22 ≤ x then (base, subtractWords lgammaOne x, 0)
        else if 0x3fc88766f1c02f69 ≤ x then (base, subtractWords x 0x3fdd8b618d5af8fc, 1)
        else (base, x, 2)
    else
      some (if 0x3ffbb4c3e8757691 ≤ x then (0, subtractWords lgammaTwo x, 0)
        else if 0x3ff310ecde3805ed ≤ x then (0, subtractWords x 0x3ff762d86356be3f, 1)
        else (0, subtractWords x lgammaOne, 2))
  return addWords base (lgammaLowCorrection arm64 fused which y)

private def lgammaMedium (arm64 fused : Bool) (x : Nat) : Option Nat := do
  let integer := magnitudeUnits x / 2 ^ 1074
  let integerWord := roundSigned (Int.ofNat integer * 2 ^ 1074) 1
  let y := subtractWords x integerWord
  let p := multiplyWords y (lgammaPolynomial arm64 y 0x3f00bfecdd17e945 [
    0x3f5e26b67368f239, 0x3f9b481c7e939961, 0x3fc2bb9cbee5f2f7, 0x3fd4d98f4f139f59,
    0x3fcb848b36e20878, 0xbfb3c467e37db0c8])
  let q := multiplyAddWords fused y (lgammaPolynomial arm64 y 0x3edebaf7a5b38140 [
    0x3f497ddaca41a95b, 0x3f9317ea742ed475, 0x3fc601edccfbdf27, 0x3fe71a1893d3dcdc,
    0x3ff645a762c4ab74]) lgammaOne
  let base := multiplyAddWords fused lgammaHalf y (divideWords p q)
  if integer ≤ 2 then return base
  let product := ((List.range (integer - 2)).reverse).foldl (fun acc k =>
    let term := roundSigned (Int.ofNat (k + 2) * 2 ^ 1074) 1
    multiplyWords acc (addWords y term)) lgammaOne
  return addWords base (← logWord arm64 (BitVec.ofNat 64 product))

private def lgammaPositive (arm64 fused : Bool) (x : Nat) : Option Nat :=
  if x = lgammaOne ∨ x = lgammaTwo then some 0
  else if x < lgammaTwo then lgammaLow arm64 fused x
  else if x < 0x4020000000000000 then lgammaMedium arm64 fused x
  else do
    let logarithm ← logWord arm64 (BitVec.ofNat 64 x)
    let reduction := subtractWords logarithm lgammaOne
    if x < 0x4390000000000000 then
      let reciprocal := divideWords lgammaOne x
      let square := multiplyWords reciprocal reciprocal
      let p := lgammaPolynomial arm64 square 0xbf5ab89d0b9e43e4 [
    0x3f4b67ba4cdad5d1, 0xbf4380cb8c0fe741, 0x3f4a019f98cf38b6, 0xbf66c16c16b02e5c,
    0x3fb555555555553b]
      let correction := multiplyAddWords arm64 reciprocal p 0x3fdacfe390c97d69
      return multiplyAddWords fused (subtractWords x lgammaHalf) reduction correction
    else return multiplyWords x reduction

/-- Go 1.27.0 src/math/lgamma.go:175–322. The first component observes
    the logarithm by word or NaN class. The second is the Gamma sign. -/
def lgammaWord (arm64 compileFused : Bool) (word : BitVec 64) : Option Nat × Int :=
  let x := word.toNat
  let magnitude := x % 2 ^ 63
  let negative := 2 ^ 63 ≤ x
  let fused := arm64 || compileFused
  if 0x7ff0000000000000 < magnitude then (none, 1)
  else if magnitude = 0x7ff0000000000000 then (some x, 1)
  else if magnitude = 0 then (some 0x7ff0000000000000, 1)
  else if magnitude < 0x3b90000000000000 then
    ((logWord arm64 (BitVec.ofNat 64 magnitude)).map negateWord, if negative then -1 else 1)
  else if !decide negative then (lgammaPositive arm64 fused magnitude, 1)
  else if 0x4330000000000000 ≤ magnitude then (some 0x7ff0000000000000, 1)
  else match sinPiWord arm64 compileFused (BitVec.ofNat 64 magnitude) with
    | none => (none, 1)
    | some sine =>
      if sine % 2 ^ 63 = 0 then (some 0x7ff0000000000000, 1)
      else
        let sign := if 2 ^ 63 ≤ sine then -1 else 1
        let divisor := multiplyWords sine magnitude % 2 ^ 63
        let value := do
          let adjustment ← logWord arm64 (BitVec.ofNat 64 (divideWords 0x400921fb54442d18 divisor))
          let positive ← lgammaPositive arm64 fused magnitude
          return subtractWords adjustment positive
        (value, sign)

theorem lgammaWord_nan (arm64 fused : Bool) (word : BitVec 64)
    (nan : 0x7ff0000000000000 < word.toNat % 2 ^ 63) :
    lgammaWord arm64 fused word = (none, 1) := by
  simp [lgammaWord, nan]

theorem lgammaWord_zero (arm64 fused : Bool) (word : BitVec 64)
    (zero : word.toNat % 2 ^ 63 = 0) :
    lgammaWord arm64 fused word = (some 0x7ff0000000000000, 1) := by
  simp [lgammaWord, zero]

theorem lgammaWord_largeNegative (arm64 fused : Bool) (word : BitVec 64)
    (negative : 2 ^ 63 ≤ word.toNat)
    (large : 0x4330000000000000 ≤ word.toNat % 2 ^ 63)
    (finite : word.toNat % 2 ^ 63 < 0x7ff0000000000000) :
    lgammaWord arm64 fused word = (some 0x7ff0000000000000, 1) := by
  have nonzero : word.toNat % 2 ^ 63 ≠ 0 := by omega
  have notTiny : ¬ word.toNat % 2 ^ 63 < 0x3b90000000000000 := by omega
  simp [lgammaWord, Nat.not_lt.mpr (Nat.le_of_lt finite), Nat.ne_of_lt finite,
    nonzero, notTiny, negative, large]

end GoLean.Math
