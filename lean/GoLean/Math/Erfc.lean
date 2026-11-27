import GoLean.Math.Exp

namespace GoLean.Math

private def erfcHorner (fused : Bool) (x leading : Nat) (remaining : List Nat) : Nat :=
  remaining.foldl (fun acc c => multiplyAddWords fused x acc c) leading

private def erfcNearZero (fused : Bool) (x : Nat) : Nat :=
  if x < 0x3c70000000000000 then x
  else
    let z := multiplyWords x x
    let r := erfcHorner fused z 0xbef8ead6120016ac [0xbf77a291236668e4, 0xbf9d2a51dbd7194f, 0xbfd4cd7d691cb913, 0x3fc06eba8214db68]
    let s := erfcHorner fused z 0xbed09c4342a26120 [0x3f215dc9221c1a10, 0x3f74d022c4d36b0f, 0x3fb0a54c5536ceba, 0x3fd97779cddadc09, 0x3ff0000000000000]
    let y := divideWords r s
    if x < 0x3fd0000000000000 then multiplyAddWords fused x y x
    else
      let half := 0x3fe0000000000000
      addWords half (multiplyAddWords fused x y (subtractWords x half))

private def erfcAroundOneCorrection (fused : Bool) (x : Nat) : Nat :=
  let s := subtractWords x 0x3ff0000000000000
  let p := erfcHorner fused s 0xbf61bf380a96073f [
      0x3fa22a36599795eb, 0xbfbc63983d3e28ec, 0x3fd45fca805120e4, 0xbfd7d240fbb8c3f1,
      0x3fda8d00ad92b34d, 0xbf6359b8bef77538]
  let q := erfcHorner fused s 0x3f888b545735151d [
      0x3f8bedc26b51dd1c, 0x3fc02660e763351f, 0x3fb2635cd99fe9a7, 0x3fe14af092eb6f33,
      0x3fbb3e6618eee323, 0x3ff0000000000000]
  divideWords p q

private def erfcTail (arm64 fused expFused : Bool) (x : Nat) : Nat :=
  let one := 0x3ff0000000000000
  let s := divideWords one (multiplyWords x x)
  let (r, denominator) :=
    if x < 0x4006db6db6db6db7 then
      (erfcHorner fused s 0xc023a0efc69ac25c [
      0xc054526557e4d2f2, 0xc067135cebccabb2, 0xc0644cb184282266, 0xc04f300ae4cba38d,
      0xc0251e0441b0e726, 0xbfe63416e4ba7360, 0xbf843412600d6435],
       erfcHorner fused s 0xbfaeeff2ee749a62 [
      0x401a47ef8e484a93, 0x405b28a3ee48ae2c, 0x407ad02157700314, 0x40842b1921ec2868,
      0x407b290dd58a1a71, 0x4061350c526ae721, 0x4033a6b9bd707687, 0x3ff0000000000000])
    else
      (erfcHorner fused s 0xc07e384e9bdc383f [
      0xc09004616a2e5992, 0xc083ec881375f228, 0xc064145d43c5ed98, 0xc031c209555f995a,
      0xbfe993ba70c285de, 0xbf84341239e86f4a],
       erfcHorner fused s 0xc03670e242712d62 [
      0x407da874e79fe763, 0x40a3f219cedf3be6, 0x40a8ffb7688c246a, 0x409802eb189d5118,
      0x40745cae221b9f0a, 0x403e568b261d5190, 0x3ff0000000000000])
  let z := x &&& 0xffffffff00000000
  let first := multiplyAddWords arm64 (negateWord z) z 0xbfe2000000000000
  let second := multiplyAddWords fused (subtractWords z x) (addWords z x) (divideWords r denominator)
  let e1 := expWord arm64 expFused (BitVec.ofNat 64 first)
  let e2 := expWord arm64 expFused (BitVec.ofNat 64 second)
  divideWords (multiplyWords e1 e2) x

/-- Go 1.27.0 src/math/erf.go:274–352. Retains the central subtraction
    order and both tail exponent stages. None denotes NaN by class. -/
def erfcWord (arm64 compileFused expFused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  let magnitude := x % 2 ^ 63
  let negative := 2 ^ 63 ≤ x
  let fused := arm64 || compileFused
  let one := 0x3ff0000000000000
  let two := 0x4000000000000000
  let erx := 0x3feb0ac160000000
  if 0x7ff0000000000000 < magnitude then none
  else if magnitude < 0x3feb000000000000 then
    let temp := erfcNearZero fused magnitude
    some (if negative then addWords one temp else subtractWords one temp)
  else if magnitude < 0x3ff4000000000000 then
    let correction := erfcAroundOneCorrection fused magnitude
    some (if negative then addWords (addWords one erx) correction
      else subtractWords (subtractWords one erx) correction)
  else if 0x403c000000000000 ≤ magnitude then some (if negative then two else 0)
  else if negative ∧ 0x4018000000000000 < magnitude then some two
  else
    let tail := erfcTail arm64 fused expFused magnitude
    some (if negative then subtractWords two tail else tail)

theorem erfcWord_nan (arm64 compileFused expFused : Bool) (word : BitVec 64)
    (nan : 0x7ff0000000000000 < word.toNat % 2 ^ 63) :
    erfcWord arm64 compileFused expFused word = none := by
  simp [erfcWord, nan]

theorem erfcWord_saturated (arm64 compileFused expFused : Bool) (word : BitVec 64)
    (large : 0x403c000000000000 ≤ word.toNat % 2 ^ 63)
    (notNaN : word.toNat % 2 ^ 63 ≤ 0x7ff0000000000000) :
    erfcWord arm64 compileFused expFused word =
      some (if 2 ^ 63 ≤ word.toNat then 0x4000000000000000 else 0) := by
  unfold erfcWord
  dsimp only
  rw [ite_eq_right (by omega), ite_eq_right (by omega), ite_eq_right (by omega), ite_eq_left large]

set_option exponentiation.threshold 4096 in
set_option maxRecDepth 4096 in
theorem erfcWord_zeros (arm64 compileFused expFused : Bool) :
    erfcWord arm64 compileFused expFused 0 = some 0x3ff0000000000000 ∧
      erfcWord arm64 compileFused expFused (BitVec.ofNat 64 0x8000000000000000) =
        some 0x3ff0000000000000 := by
  cases arm64 <;> cases compileFused <;> cases expFused <;> decide

end GoLean.Math
