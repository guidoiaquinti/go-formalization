import GoLean.Math.Exp

namespace GoLean.Math

private def erfHorner (fused : Bool) (x leading : Nat) (remaining : List Nat) : Nat :=
  remaining.foldl (fun acc c => multiplyAddWords fused x acc c) leading

private def erfNearZero (arm64 fused : Bool) (x : Nat) : Nat :=
  if x < 0x0080000000000000 then
    let combined := if arm64 then
      multiplyAddWords true x 0x3ff06eba8214db69 (multiplyWords 0x4020000000000000 x)
    else
      multiplyAddWords fused x 0x4020000000000000 (multiplyWords 0x3ff06eba8214db69 x)
    multiplyWords 0x3fc0000000000000 combined
  else if x < 0x3e30000000000000 then
    multiplyAddWords fused x 0x3fc06eba8214db69 x
  else
    let z := multiplyWords x x
    let r := erfHorner fused z 0xbef8ead6120016ac [0xbf77a291236668e4, 0xbf9d2a51dbd7194f, 0xbfd4cd7d691cb913, 0x3fc06eba8214db68]
    let s := erfHorner fused z 0xbed09c4342a26120 [0x3f215dc9221c1a10, 0x3f74d022c4d36b0f, 0x3fb0a54c5536ceba, 0x3fd97779cddadc09, 0x3ff0000000000000]
    multiplyAddWords fused x (divideWords r s) x

private def erfAroundOne (fused : Bool) (x : Nat) : Nat :=
  let s := subtractWords x 0x3ff0000000000000
  let p := erfHorner fused s 0xbf61bf380a96073f [
      0x3fa22a36599795eb, 0xbfbc63983d3e28ec, 0x3fd45fca805120e4, 0xbfd7d240fbb8c3f1,
      0x3fda8d00ad92b34d, 0xbf6359b8bef77538]
  let q := erfHorner fused s 0x3f888b545735151d [
      0x3f8bedc26b51dd1c, 0x3fc02660e763351f, 0x3fb2635cd99fe9a7, 0x3fe14af092eb6f33,
      0x3fbb3e6618eee323, 0x3ff0000000000000]
  addWords 0x3feb0ac160000000 (divideWords p q)

private def erfTail (arm64 fused expFused : Bool) (x : Nat) : Nat :=
  let one := 0x3ff0000000000000
  let s := divideWords one (multiplyWords x x)
  let (r, denominator) :=
    if x < 0x4006db6db6db6db7 then
      (erfHorner fused s 0xc023a0efc69ac25c [
      0xc054526557e4d2f2, 0xc067135cebccabb2, 0xc0644cb184282266, 0xc04f300ae4cba38d,
      0xc0251e0441b0e726, 0xbfe63416e4ba7360, 0xbf843412600d6435],
       erfHorner fused s 0xbfaeeff2ee749a62 [
      0x401a47ef8e484a93, 0x405b28a3ee48ae2c, 0x407ad02157700314, 0x40842b1921ec2868,
      0x407b290dd58a1a71, 0x4061350c526ae721, 0x4033a6b9bd707687, 0x3ff0000000000000])
    else
      (erfHorner fused s 0xc07e384e9bdc383f [
      0xc09004616a2e5992, 0xc083ec881375f228, 0xc064145d43c5ed98, 0xc031c209555f995a,
      0xbfe993ba70c285de, 0xbf84341239e86f4a],
       erfHorner fused s 0xc03670e242712d62 [
      0x407da874e79fe763, 0x40a3f219cedf3be6, 0x40a8ffb7688c246a, 0x409802eb189d5118,
      0x40745cae221b9f0a, 0x403e568b261d5190, 0x3ff0000000000000])
  let z := x &&& 0xffffffff00000000
  let first := multiplyAddWords arm64 (negateWord z) z 0xbfe2000000000000
  let second := multiplyAddWords fused (subtractWords z x) (addWords z x) (divideWords r denominator)
  let e1 := expWord arm64 expFused (BitVec.ofNat 64 first)
  let e2 := expWord arm64 expFused (BitVec.ofNat 64 second)
  subtractWords one (divideWords (multiplyWords e1 e2) x)

/-- Go 1.27.0 src/math/erf.go:189–265. The model follows the tiny
    scaling, central rational polynomials, truncated-word tail Exp
    arguments and sign restoration. None denotes NaN by class. -/
def erfWord (arm64 compileFused expFused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  let magnitude := x % 2 ^ 63
  let fused := arm64 || compileFused
  let signed := fun result => if 2 ^ 63 ≤ x then negateWord result else result
  if 0x7ff0000000000000 < magnitude then none
  else if magnitude = 0 then some x
  else if 0x4018000000000000 ≤ magnitude then some (signed 0x3ff0000000000000)
  else if magnitude < 0x3feb000000000000 then some (signed (erfNearZero arm64 fused magnitude))
  else if magnitude < 0x3ff4000000000000 then some (signed (erfAroundOne fused magnitude))
  else some (signed (erfTail arm64 fused expFused magnitude))

theorem erfWord_zero (arm64 compileFused expFused : Bool) (word : BitVec 64)
    (zero : word.toNat % 2 ^ 63 = 0) :
    erfWord arm64 compileFused expFused word = some word.toNat := by
  simp [erfWord, zero]

theorem erfWord_nan (arm64 compileFused expFused : Bool) (word : BitVec 64)
    (nan : 0x7ff0000000000000 < word.toNat % 2 ^ 63) :
    erfWord arm64 compileFused expFused word = none := by
  simp [erfWord, nan]

theorem erfWord_saturated (arm64 compileFused expFused : Bool) (word : BitVec 64)
    (large : 0x4018000000000000 ≤ word.toNat % 2 ^ 63)
    (notNaN : word.toNat % 2 ^ 63 ≤ 0x7ff0000000000000) :
    erfWord arm64 compileFused expFused word =
      some (if 2 ^ 63 ≤ word.toNat then 0xbff0000000000000 else 0x3ff0000000000000) := by
  have nonzero : word.toNat % 2 ^ 63 ≠ 0 := by omega
  simp [erfWord, Nat.not_lt.mpr notNaN, nonzero, large, negateWord]

end GoLean.Math
