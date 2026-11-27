import GoLean.Math.Exp2

namespace GoLean.Math

private def expReducedOperands (fuseNegative fusePositive fuseHi : Bool) (x : Nat) :
    Nat × Nat × Int :=
  let negative := 2 ^ 63 ≤ x
  let half := if negative then 0xbfe0000000000000 else 0x3fe0000000000000
  let k := (signedUnits (multiplyAddWords (if negative then fuseNegative else fusePositive)
    x 0x3ff71547652b82fe half)).tdiv (2 ^ 1074)
  let kWord := roundSigned (k * 2 ^ 1074) 1
  let hi := multiplyAddWords fuseHi kWord 0xbfe62e42fee00000 x
  let lo := multiplyWords kWord 0x3dea39ef35793c76
  (hi, lo, k)

private def expArm64 (x : Nat) : Nat :=
  let (hi, lo, k) := expReducedOperands true true true x
  scaleExp2 true (exp2Reduced true true hi lo) k

private def roundExpIndex (word : Nat) : Int :=
  let units := signedUnits word
  let position := units.natAbs / 2 ^ 1074
  let rem := units.natAbs % 2 ^ 1074
  let rounded : Nat := if 2 ^ 1073 < rem ∨ (rem = 2 ^ 1073 ∧ position % 2 = 1) then position + 1 else position
  if units < 0 then -(rounded : Int) else (rounded : Int)

private def expAmd64 (fused : Bool) (x : Nat) : Nat :=
  let k := roundExpIndex (multiplyWords x 0x3ff71547652b82fe)
  if k < -1075 then 0
  else if 1023 < k then 0x7ff0000000000000
  else
    let kWord := roundSigned (k * 2 ^ 1074) 1
    let upper := multiplyAddWords fused kWord 0xbfe62e42fefa3000 x
    let lower := multiplyAddWords fused kWord 0xbd53de6af278ece6 upper
    let r := multiplyWords lower 0x3fb0000000000000
    let term := multiplyAddWords fused r 0x3efa01a01a01a01a 0x3f2a01a01a01a01a
    let term := multiplyAddWords fused r term 0x3f56c16c16c16c17
    let term := multiplyAddWords fused r term 0x3f81111111111111
    let term := multiplyAddWords fused r term 0x3fa5555555555555
    let term := multiplyAddWords fused r term 0x3fc5555555555555
    let term := multiplyAddWords fused r term 0x3fe0000000000000
    let term := multiplyAddWords fused r term 0x3ff0000000000000
    let term := multiplyWords r term
    let double := fun t => multiplyWords t (addWords 0x4000000000000000 t)
    let term := double (double (double term))
    let result := multiplyAddWords fused term (addWords 0x4000000000000000 term) 0x3ff0000000000000
    let biased := k + 1023
    if 0 < biased then multiplyWords result (biased.toNat * 2 ^ 52)
    else multiplyWords (multiplyWords result ((biased + 1022).toNat * 2 ^ 52)) (2 ^ 52)

/-- Go 1.27.0 src/math/exp.go:16–21, exp_arm64.s:26–105 and
    exp_amd64.s:40–157. Arm64 uses a split reduction, rational correction
    and exponent-word scale. Amd64 uses its own rounded reduction,
    Taylor evaluation, four doubling updates and staged scale.
    The amd64 fused flag selects the AVX/FMA instructions. -/
def expWord (arm64 fused : Bool) (word : BitVec 64) : Nat :=
  let x := word.toNat
  let magnitude := x % 2 ^ 63
  if 0x7ff0000000000000 < magnitude then x
  else if x < 2 ^ 63 ∧ 0x40862e42fefa39ef < magnitude then 0x7ff0000000000000
  else if 2 ^ 63 ≤ x ∧ (if arm64 then 0x40874910d52d3051 else 0x4087500000000000) < magnitude then 0
  else if arm64 then
    if magnitude < 0x3e30000000000000 then addWords 0x3ff0000000000000 x
    else expArm64 x
  else expAmd64 fused x

theorem expWord_nan (arm64 fused : Bool) (word : BitVec 64)
    (nan : 0x7ff0000000000000 < word.toNat % 2 ^ 63) :
    expWord arm64 fused word = word.toNat := by
  simp [expWord, nan]

theorem expWord_infinities (arm64 fused : Bool) :
    expWord arm64 fused (BitVec.ofNat 64 0x7ff0000000000000) = 0x7ff0000000000000 ∧
    expWord arm64 fused (BitVec.ofNat 64 0xfff0000000000000) = 0 := by
  cases arm64 <;> cases fused <;> decide

theorem expWord_nearZero (fused : Bool) (word : BitVec 64)
    (near : word.toNat % 2 ^ 63 < 0x3e30000000000000) :
    expWord true fused word = addWords 0x3ff0000000000000 word.toNat := by
  unfold expWord
  dsimp only
  rw [ite_eq_right (by omega), ite_eq_right (by omega)]
  simp [near, show ¬ 0x40874910d52d3051 < word.toNat % 2 ^ 63 by omega]

/-- Go 1.27.0 src/math/exp.go:99–134. The generic source shares the
    fdlibm polynomial and correction with arm64 assembly and uses Go
    Ldexp scaling. The pinned amd64 v3/v4 compiler contracts the positive
    half-offset and polynomial additions. Arm64 also contracts the
    negative half-offset and both multiply-subtract operations. -/
def expGenericWord (arm64 fused : Bool) (word : BitVec 64) : Nat :=
  let x := word.toNat
  let magnitude := x % 2 ^ 63
  if 0x7ff0000000000000 < magnitude then x
  else if x < 2 ^ 63 ∧ 0x40862e42fefa39ef < magnitude then 0x7ff0000000000000
  else if 2 ^ 63 ≤ x ∧ 0x40874910d52d3051 < magnitude then 0
  else if magnitude < 0x3e30000000000000 then addWords 0x3ff0000000000000 x
  else
    let (hi, lo, k) := expReducedOperands arm64 (arm64 || fused) arm64 x
    ldexp 64 (exp2Reduced (arm64 || fused) arm64 hi lo) k

theorem expGenericWord_nan (arm64 fused : Bool) (word : BitVec 64)
    (nan : 0x7ff0000000000000 < word.toNat % 2 ^ 63) :
    expGenericWord arm64 fused word = word.toNat := by
  simp [expGenericWord, nan]

theorem expGenericWord_infinities (arm64 fused : Bool) :
    expGenericWord arm64 fused (BitVec.ofNat 64 0x7ff0000000000000) = 0x7ff0000000000000 ∧
    expGenericWord arm64 fused (BitVec.ofNat 64 0xfff0000000000000) = 0 := by
  cases arm64 <;> cases fused <;> decide

theorem expGenericWord_nearZero (arm64 fused : Bool) (word : BitVec 64)
    (near : word.toNat % 2 ^ 63 < 0x3e30000000000000) :
    expGenericWord arm64 fused word = addWords 0x3ff0000000000000 word.toNat := by
  unfold expGenericWord
  dsimp only
  rw [ite_eq_right (by omega), ite_eq_right (by omega), ite_eq_right (by omega), ite_eq_left near]

end GoLean.Math
