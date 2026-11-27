import GoLean.Math.Exp2
import GoLean.Math.Log

namespace GoLean.Math

private def erfinvHorner (fused : Bool) (x leading : Nat) (remaining : List Nat) : Nat :=
  remaining.foldl (fun acc c => multiplyAddWords fused x acc c) leading

private def erfinvCentralRatio (fused : Bool) (r : Nat) : Nat × Nat :=
  let numerator := erfinvHorner fused r 0x408bb8c0a7936704 [
    0x40c715bf25ff1d11, 0x40d73982a6012afb, 0x40cfb5ee66e5a285, 0x40b2f6e30b2ec51e,
    0x4085c88056b01502, 0x4047894dd814bfa3, 0x3ff32917a42157d1]
  let denominator := erfinvHorner fused r 0x40b46a7eca984b69 [
    0x40dc0e457cb1ae76, 0x40e3317caa64f4be, 0x40d4b772d5d65266, 0x40b512322e75c89f,
    0x4085797efdc8b3f7, 0x4045281b386e1ab5, 0x3ff0000000000000]
  (numerator, denominator)

private def erfinvFirstTailRatio (fused : Bool) (r : Nat) : Nat × Nat :=
  let numerator := erfinvHorner fused r 0x3f49615ac0b7ace9 [
    0x3f9744eb6c45ec67, 0x3fcef2abb9b85c37, 0x3ff453cc085375b2, 0x400d2ecb1a3d02c4,
    0x401713f71462256a, 0x4012857748cab19b, 0x3ff6c665fde9526a]
  let denominator := erfinvHorner fused r 0x3e19876e6013e192 [
    0x3f496042ab9205be, 0x3f960290af9f12bc, 0x3fcacf476a756d3d, 0x3fef371e4f4de1a0,
    0x4002f7543ff6ba47, 0x40073aad9bca5405, 0x3ff6a09e667f3bcd]
  (numerator, denominator)

private def erfinvSecondTailRatio (fused : Bool) (r : Nat) : Nat × Nat :=
  let numerator := erfinvHorner fused r 0x3e8afb74d693bf93 [
    0x3efc6ec6cc59e02a, 0x3f545c1908425345, 0x3f9b2b41193b4ee7, 0x3fd2fad9315255cf,
    0x3ffc8ea6461fa445, 0x4015daea6e875003, 0x401aa1b1c13ee526]
  let denominator := erfinvHorner fused r 0x3cea0a3fb665cb0a [
    0x3e8afb693c174d20, 0x3efb61164a01861a, 0x3f523b6cf8ec177a, 0x3f958ab4fd6bcd50,
    0x3fc8c975f4d2f24b, 0x3feb2532d5cd90e5, 0x3ff6a09e667f3bcd]
  (numerator, denominator)

/-- Go 1.27.0 src/math/erfinv.go:77–116. Retains the rounded central
    radicand, the Log and Sqrt tail reduction and all three rational
    approximations. None denotes NaN by class. -/
def erfinvWord (arm64 compileFused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  let magnitude := x % 2 ^ 63
  let fused := arm64 || compileFused
  let signed := fun result => if 2 ^ 63 ≤ x then negateWord result else result
  if 0x3ff0000000000000 < magnitude then none
  else if magnitude = 0x3ff0000000000000 then some (signed 0x7ff0000000000000)
  else if magnitude = 0 then some x
  else if magnitude ≤ 0x3feb333333333333 then
    let scaled := multiplyWords 0x3fd0000000000000 magnitude
    let r := multiplyAddWords arm64 scaled (negateWord magnitude) 0x3fc71eb851eb851f
    let (p, q) := erfinvCentralRatio fused r
    some (signed (divideWords (multiplyWords magnitude p) q))
  else do
    let difference := subtractWords 0x3ff0000000000000 magnitude
    let logarithm ← logWord arm64 (BitVec.ofNat 64 difference)
    let r ← sqrtWord (BitVec.ofNat 64 (subtractWords 0x3fe62e42fefa39ef logarithm))
    let (p, q) := if r ≤ 0x4014000000000000 then
        erfinvFirstTailRatio fused (subtractWords r 0x3ff999999999999a)
      else
        erfinvSecondTailRatio fused (subtractWords r 0x4014000000000000)
    return signed (divideWords p q)

theorem erfinvWord_zero (arm64 fused : Bool) (word : BitVec 64)
    (zero : word.toNat % 2 ^ 63 = 0) :
    erfinvWord arm64 fused word = some word.toNat := by
  simp [erfinvWord, zero]

theorem erfinvWord_outside (arm64 fused : Bool) (word : BitVec 64)
    (outside : 0x3ff0000000000000 < word.toNat % 2 ^ 63) :
    erfinvWord arm64 fused word = none := by
  simp [erfinvWord, outside]

theorem erfinvWord_endpoints (arm64 fused : Bool) :
    erfinvWord arm64 fused (BitVec.ofNat 64 0x3ff0000000000000) = some 0x7ff0000000000000 ∧
      erfinvWord arm64 fused (BitVec.ofNat 64 0xbff0000000000000) = some 0xfff0000000000000 := by
  cases arm64 <;> cases fused <;> decide

end GoLean.Math
