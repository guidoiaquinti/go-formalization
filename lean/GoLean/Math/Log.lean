import GoLean.Math.Exp2
import GoLean.Math.Frexp

namespace GoLean.Math

private def logReduction (normalize : Bool) (word : Nat) : Nat × Int :=
  let (fraction, exponent) :=
    if normalize then frexp word
    else (0x3fe0000000000000 + word % 2 ^ 52, ((word >>> 52 : Nat) : Int) - 1022)
  if fraction < 0x3fe6a09e667f3bcd then
    (multiplyWords fraction 0x4000000000000000, exponent - 1)
  else (fraction, exponent)

private def logEvaluate (arm64 fused : Bool) (fraction : Nat) (exponent : Int) : Nat :=
  let k := roundSigned (exponent * 2 ^ 1074) 1
  let f := subtractWords fraction 0x3ff0000000000000
  let s := divideWords f (addWords 0x4000000000000000 f)
  let s2 := multiplyWords s s
  let s4 := multiplyWords s2 s2
  let odd := multiplyAddWords fused s4 0x3fc2f112df3e5244 0x3fc7466496cb03de
  let odd := multiplyAddWords fused s4 odd 0x3fd2492494229359
  let odd := multiplyAddWords fused s4 odd 0x3fe5555555555593
  let even := multiplyAddWords fused s4 0x3fc39a09d078c69f 0x3fcc71c51d8e78af
  let even := multiplyAddWords fused s4 even 0x3fd999999997fa04
  let poly := multiplyAddWords fused s2 odd (multiplyWords s4 even)
  let half := multiplyWords 0x3fe0000000000000 f
  let sum := multiplyAddWords fused f half poly
  let correction := if arm64 then multiplyAddWords true s sum (multiplyWords k 0x3dea39ef35793c76)
    else multiplyAddWords fused k 0x3dea39ef35793c76 (multiplyWords s sum)
  let residual := multiplyAddWords arm64 f half (negateWord correction)
  multiplyAddWords arm64 k 0x3fe62e42fee00000 (negateWord (subtractWords residual f))

private def logWithProfile (normalize arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  let magnitude := x % 2 ^ 63
  if magnitude = 0 then some 0xfff0000000000000
  else if 2 ^ 63 ≤ x ∨ 0x7ff0000000000000 < magnitude then none
  else if x = 0x7ff0000000000000 then some x
  else
    let (fraction, exponent) := logReduction normalize x
    some (logEvaluate arm64 fused fraction exponent)

/-- Go 1.27.0 src/math/log.go:81–129 and log_amd64.s:22–112.
    The arm64 compiler contracts selected multiply-add operations and
    normalizes denormals through Frexp. The amd64 assembly uses a forced
    leading bit even when the exponent field is zero. -/
def logWord (arm64 : Bool) (word : BitVec 64) : Option Nat :=
  logWithProfile arm64 arm64 arm64 word

/-- Go 1.27.0 src/math/log.go:88–129. The generic helper always
    normalizes subnormals through Frexp. Amd64 v3/v4 contracts selected
    additions but rounds the final subtractions separately. -/
def logGenericWord (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  logWithProfile true arm64 (arm64 || fused) word

theorem logWord_zero (arm64 : Bool) (word : BitVec 64)
    (zero : word.toNat % 2 ^ 63 = 0) :
    logWord arm64 word = some 0xfff0000000000000 := by
  simp [logWord, logWithProfile, zero]

theorem logWord_negative (arm64 : Bool) (word : BitVec 64)
    (negative : 2 ^ 63 ≤ word.toNat) (nonzero : word.toNat % 2 ^ 63 ≠ 0) :
    logWord arm64 word = none := by
  simp [logWord, logWithProfile, nonzero, negative]

theorem logWord_inf (arm64 : Bool) :
    logWord arm64 (BitVec.ofNat 64 0x7ff0000000000000) = some 0x7ff0000000000000 := by
  cases arm64 <;> decide

theorem logGenericWord_zero (arm64 fused : Bool) (word : BitVec 64)
    (zero : word.toNat % 2 ^ 63 = 0) :
    logGenericWord arm64 fused word = some 0xfff0000000000000 := by
  simp [logGenericWord, logWithProfile, zero]

theorem logGenericWord_negative (arm64 fused : Bool) (word : BitVec 64)
    (negative : 2 ^ 63 ≤ word.toNat) (nonzero : word.toNat % 2 ^ 63 ≠ 0) :
    logGenericWord arm64 fused word = none := by
  simp [logGenericWord, logWithProfile, nonzero, negative]

theorem logGenericWord_inf (arm64 fused : Bool) :
    logGenericWord arm64 fused (BitVec.ofNat 64 0x7ff0000000000000) = some 0x7ff0000000000000 := by
  cases arm64 <;> cases fused <;> decide

end GoLean.Math
