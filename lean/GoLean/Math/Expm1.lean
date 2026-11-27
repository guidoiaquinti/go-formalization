import GoLean.Math.Exp2

namespace GoLean.Math

private def expm1Reduction (arm64 fused : Bool) (raw magnitude : Nat) : Nat × Nat × Int :=
  let negative := 2 ^ 63 ≤ raw
  let (hi, lo, k) :=
    if magnitude < 0x3ff0a2b23f3bab73 then
      if negative then (addWords raw 0x3fe62e42fee00000, 0xbdea39ef35793c76, (-1 : Int))
      else (subtractWords raw 0x3fe62e42fee00000, 0x3dea39ef35793c76, (1 : Int))
    else
      let offset := if negative then 0xbfe0000000000000 else 0x3fe0000000000000
      let contract := arm64 || (fused && !decide negative)
      let k := (signedUnits (multiplyAddWords contract raw 0x3ff71547652b82fe offset)).tdiv (2 ^ 1074)
      let integer := roundSigned (k * 2 ^ 1074) 1
      (multiplyAddWords arm64 integer 0xbfe62e42fee00000 raw, multiplyWords integer 0x3dea39ef35793c76, k)
  let reduced := subtractWords hi lo
  (reduced, subtractWords (subtractWords hi reduced) lo, k)

private def addExpField (x : Nat) (k : Int) : Nat :=
  ((x : Int) + k * 2 ^ 52).toNat

private def expm1Polynomial (arm64 fused : Bool) (x correction : Nat) (k : Int) : Nat :=
  let half := multiplyWords 0x3fe0000000000000 x
  let square := multiplyWords x half
  let polynomial := multiplyAddWords fused square 0xbe8afdb76e09c32d 0x3ed0cfca86e65239
  let polynomial := multiplyAddWords fused square polynomial 0xbf14ce199eaadbb7
  let polynomial := multiplyAddWords fused square polynomial 0x3f5a01a019fe5585
  let polynomial := multiplyAddWords fused square polynomial 0xbfa11111111110f4
  let r1 := multiplyAddWords fused square polynomial 0x3ff0000000000000
  let t := multiplyAddWords arm64 r1 (negateWord half) 0x4008000000000000
  let divisor := multiplyAddWords arm64 x (negateWord t) 0x4018000000000000
  let ratio := divideWords (subtractWords r1 t) divisor
  -- On arm64 the compiler retains x*e as a rounded product here and
  -- contracts the subtraction of x*hfx into that product.
  if k = 0 then
    let error := multiplyAddWords arm64 x (negateWord half) (multiplyWords x (multiplyWords square ratio))
    subtractWords x error
  else
    -- The arm64 compiler contracts both corrections and the subtraction
    -- of x*hfx. Amd64 rounds each of those stages separately.
    let adjusted := multiplyAddWords arm64 square ratio (negateWord correction)
    let tail := multiplyAddWords arm64 x adjusted (negateWord correction)
    let error := multiplyAddWords arm64 x (negateWord half) tail
    if k = -1 then multiplyAddWords arm64 0x3fe0000000000000 (subtractWords x error) 0xbfe0000000000000
    else if k = 1 then
      if 0xbfd0000000000000 < x then
        multiplyWords 0xc000000000000000 (subtractWords error (addWords x 0x3fe0000000000000))
      else addWords 0x3ff0000000000000 (multiplyWords 0x4000000000000000 (subtractWords x error))
    else if k ≤ -2 ∨ 56 < k then
      let scaled := addExpField (subtractWords 0x3ff0000000000000 (subtractWords error x)) k
      if 0x7ff0000000000000 ≤ scaled then 0x7ff0000000000000
      else subtractWords scaled 0x3ff0000000000000
    else if k < 20 then
      let t := 0x3ff0000000000000 - (0x20000000000000 >>> k.toNat)
      addExpField (subtractWords t (subtractWords error x)) k
    else
      let t := (1023 - k).toNat <<< 52
      addExpField (addWords (subtractWords x (addWords error t)) 0x3ff0000000000000) k

/-- Go 1.27.0 src/math/expm1.go:126–244. Preserve the argument reduction,
    correction polynomial, and exponent-word restoration branches with
    the rounding and contractions of the pinned compiler profile. -/
def expm1Word (arm64 fused : Bool) (word : BitVec 64) : Nat :=
  let x := word.toNat
  let magnitude := x % 2 ^ 63
  if 0x7ff0000000000000 ≤ magnitude then
    if x = 0xfff0000000000000 then 0xbff0000000000000 else x
  else if 2 ^ 63 ≤ x ∧ 0x4043687a9f1af2b1 ≤ magnitude then 0xbff0000000000000
  else if x < 2 ^ 63 ∧ 0x40862e42fefa39ef ≤ magnitude then 0x7ff0000000000000
  else if 0x3fd62e42fefa39ef < magnitude then
    let (reduced, correction, k) := expm1Reduction arm64 fused x magnitude
    expm1Polynomial arm64 (arm64 || fused) reduced correction k
  else if magnitude < 0x3c90000000000000 then x
  else expm1Polynomial arm64 (arm64 || fused) x 0 0

theorem expm1Word_nan (arm64 fused : Bool) (word : BitVec 64)
    (nan : 0x7ff0000000000000 < word.toNat % 2 ^ 63) :
    expm1Word arm64 fused word = word.toNat := by
  unfold expm1Word
  dsimp only
  rw [ite_eq_left (by omega), ite_eq_right (by omega)]

theorem expm1Word_tiny (arm64 fused : Bool) (word : BitVec 64)
    (tiny : word.toNat % 2 ^ 63 < 0x3c90000000000000) :
    expm1Word arm64 fused word = word.toNat := by
  unfold expm1Word
  dsimp only
  rw [ite_eq_right (by omega), ite_eq_right (by omega), ite_eq_right (by omega),
    ite_eq_right (by omega), ite_eq_left tiny]

theorem expm1Word_infinities (arm64 fused : Bool) :
    expm1Word arm64 fused (BitVec.ofNat 64 0x7ff0000000000000) = 0x7ff0000000000000 ∧
    expm1Word arm64 fused (BitVec.ofNat 64 0xfff0000000000000) = 0xbff0000000000000 := by
  cases arm64 <;> cases fused <;> decide

end GoLean.Math
