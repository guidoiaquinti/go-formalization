import GoLean.Math.Exp2

namespace GoLean.Math

private def log1pReduce (x magnitude : Nat) : Nat × Nat × Int × Nat :=
  -- These encodings come from the pinned compiler's decimal constants.
  if magnitude < 0x3fda827999fcef32 ∧ x < 0xbfd2bec333018867 then
    (x, 0, 0, 1)
  else
    let roundInput := magnitude < 0x4340000000000000
    let u := if roundInput then addWords 0x3ff0000000000000 x else x
    let k := ((u >>> 52 : Nat) : Int) - 1023
    let correction :=
      if roundInput then
        let difference := if 0 < k then subtractWords 0x3ff0000000000000 (subtractWords u x)
          else subtractWords x (subtractWords u 0x3ff0000000000000)
        divideWords difference u
      else 0
    let fraction := u % 2 ^ 52
    if fraction < 0x6a09e667f3bcd then
      (subtractWords (0x3ff0000000000000 + fraction) 0x3ff0000000000000, correction, k, fraction)
    else
      (subtractWords (0x3fe0000000000000 + fraction) 0x3ff0000000000000, correction, k + 1,
        (2 ^ 52 - fraction) >>> 2)

private def log1pSeries (fused : Bool) (z : Nat) : Nat :=
  let p := multiplyAddWords fused z 0x3fc2f112df3e5244 0x3fc39a09d078c69f
  let p := multiplyAddWords fused z p 0x3fc7466496cb03de
  let p := multiplyAddWords fused z p 0x3fcc71c51d8e78af
  let p := multiplyAddWords fused z p 0x3fd2492494229359
  let p := multiplyAddWords fused z p 0x3fd999999997fa04
  multiplyAddWords fused z p 0x3fe5555555555593

private def log1pEvaluate (arm64 fused : Bool) (f correction : Nat) (k : Int) (fraction : Nat) : Nat :=
  let fused := arm64 || fused
  let half := multiplyWords 0x3fe0000000000000 f
  let square := multiplyWords half f
  let integer := roundSigned (k * 2 ^ 1074) 1
  let tail := multiplyAddWords fused integer 0x3dea39ef35793c76 correction
  let join := fun body => multiplyAddWords arm64 integer 0x3fe62e42fee00000 (negateWord body)
  if fraction = 0 then
    if f % 2 ^ 63 = 0 then
      if k = 0 then 0 else multiplyAddWords fused integer 0x3fe62e42fee00000 tail
    else
      let factor := multiplyAddWords arm64 f 0xbfe5555555555555 0x3ff0000000000000
      if k = 0 then multiplyAddWords arm64 square (negateWord factor) f
      else join (subtractWords (multiplyAddWords arm64 square factor (negateWord tail)) f)
  else
    let s := divideWords f (addWords 0x4000000000000000 f)
    let z := multiplyWords s s
    let poly := log1pSeries fused z
    -- Arm64 contracts f*hfx + R. Amd64 v3/v4 contracts z*poly + hfsq.
    let correctionSum := if arm64 then multiplyAddWords true f half (multiplyWords z poly)
      else multiplyAddWords fused z poly square
    if k = 0 then
      subtractWords f (multiplyAddWords arm64 s (negateWord correctionSum) square)
    else
      let correction := multiplyAddWords fused s correctionSum tail
      let residual := multiplyAddWords arm64 f half (negateWord correction)
      join (subtractWords residual f)

/-- Go 1.27.0 src/math/log1p.go:96–203. Retains the compensated
    normalization, tiny-input approximations, interval polynomial, and
    final split-ln(2) reconstruction. NaN results are checked by class. -/
def log1pWord (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  let magnitude := x % 2 ^ 63
  if 0x7ff0000000000000 < magnitude ∨ 0xbff0000000000000 < x then none
  else if x = 0xbff0000000000000 then some 0xfff0000000000000
  else if x = 0x7ff0000000000000 then some x
  else if magnitude < 0x3c90000000000000 then some x
  else if magnitude < 0x3e20000000000000 then
    some (multiplyAddWords arm64 (multiplyWords x x) 0xbfe0000000000000 x)
  else
    let (f, correction, k, fraction) := log1pReduce x magnitude
    some (log1pEvaluate arm64 fused f correction k fraction)

theorem log1pWord_below (arm64 fused : Bool) (word : BitVec 64)
    (below : 0xbff0000000000000 < word.toNat) : log1pWord arm64 fused word = none := by
  simp [log1pWord, below]

theorem log1pWord_tiny (arm64 fused : Bool) (word : BitVec 64)
    (tiny : word.toNat % 2 ^ 63 < 0x3c90000000000000) :
    log1pWord arm64 fused word = some word.toNat := by
  have := word.isLt
  unfold log1pWord
  dsimp only
  rw [ite_eq_right (by omega), ite_eq_right (by omega), ite_eq_right (by omega), ite_eq_left tiny]

theorem log1pWord_edges (arm64 fused : Bool) :
    log1pWord arm64 fused (BitVec.ofNat 64 0xbff0000000000000) = some 0xfff0000000000000 ∧
    log1pWord arm64 fused (BitVec.ofNat 64 0x7ff0000000000000) = some 0x7ff0000000000000 := by
  cases arm64 <;> cases fused <;> decide

end GoLean.Math
