import GoLean.Math.Exp

namespace GoLean.Math

private def tanhRational (arm64 compileFused : Bool) (x : Nat) : Nat :=
  let s := multiplyWords x x
  let p := multiplyAddWords arm64 0xbfeedc5baafd6f4b s 0xc058d26a0e26682d
  let p := multiplyAddWords arm64 p s 0xc0993ac030580563
  let q := multiplyAddWords (arm64 || compileFused) x x 0x405c33f28a581b86
  let q := multiplyAddWords arm64 q s 0x40a176fa0e5535fa
  let q := multiplyAddWords arm64 q s 0x40b2ec102442040c
  addWords x (divideWords (multiplyWords (multiplyWords x s) p) q)

/-- Go 1.27.0 src/math/tanh.go:82–105. The two amd64 flags select
    compiler contraction in the rational branch and runtime fusion in Exp.
    None represents NaN by class. -/
def tanhWord (arm64 compileFused expFused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  let magnitude := x % 2 ^ 63
  let signed := fun result => if 2 ^ 63 ≤ x then negateWord result else result
  if 0x7ff0000000000000 < magnitude then none
  else if 0x404601e678fc457b < magnitude then some (signed 0x3ff0000000000000)
  else if 0x3fe4000000000000 ≤ magnitude then
    let doubled := multiplyWords 0x4000000000000000 magnitude
    let s := expWord arm64 expFused (BitVec.ofNat 64 doubled)
    let q := divideWords 0x4000000000000000 (addWords s 0x3ff0000000000000)
    some (signed (subtractWords 0x3ff0000000000000 q))
  else if magnitude = 0 then some x
  else some (tanhRational arm64 compileFused x)

theorem tanhWord_zero (arm64 compileFused expFused : Bool) (word : BitVec 64)
    (zero : word.toNat % 2 ^ 63 = 0) :
    tanhWord arm64 compileFused expFused word = some word.toNat := by
  simp [tanhWord, zero]

theorem tanhWord_nan (arm64 compileFused expFused : Bool) (word : BitVec 64)
    (nan : 0x7ff0000000000000 < word.toNat % 2 ^ 63) :
    tanhWord arm64 compileFused expFused word = none := by
  simp [tanhWord, nan]

theorem tanhWord_saturated (arm64 compileFused expFused : Bool) (word : BitVec 64)
    (large : 0x404601e678fc457b < word.toNat % 2 ^ 63)
    (notNaN : word.toNat % 2 ^ 63 ≤ 0x7ff0000000000000) :
    tanhWord arm64 compileFused expFused word =
      some (if 2 ^ 63 ≤ word.toNat then 0xbff0000000000000 else 0x3ff0000000000000) := by
  simp [tanhWord, Nat.not_lt.mpr notNaN, large, negateWord]

end GoLean.Math
