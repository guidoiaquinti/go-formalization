import Init.Data.Int.DivMod.Lemmas
import Init.Omega

namespace GoLean.Math.Big

/-- Value-level model of Go 1.27.0 src/math/big/int.go:273–305.
    An absent result denotes the division-by-zero panic. -/
def intQuoRem (x y : Int) : Option (Int × Int) :=
  if y = 0 then none else some (x.tdiv y, x.tmod y)

def intQuo (x y : Int) : Option Int := (intQuoRem x y).map Prod.fst

def intRem (x y : Int) : Option Int := (intQuoRem x y).map Prod.snd

/-- The pinned DivMod adjusts the truncated pair when the remainder is negative.
    The original divisor remains available for the correction. -/
def intDivMod (x y : Int) : Option (Int × Int) :=
  (intQuoRem x y).map fun (q, r) =>
    if r < 0 then
      if y < 0 then (q + 1, r - y) else (q - 1, r + y)
    else (q, r)

def intDiv (x y : Int) : Option Int := (intDivMod x y).map Prod.fst

def intMod (x y : Int) : Option Int := (intDivMod x y).map Prod.snd

theorem intQuo_truncates (x y : Int) (nonzero : y ≠ 0) :
    intQuo x y = some (x.tdiv y) := by simp [intQuo, intQuoRem, nonzero]

theorem intRem_truncates (x y : Int) (nonzero : y ≠ 0) :
    intRem x y = some (x.tmod y) := by simp [intRem, intQuoRem, nonzero]

theorem intQuoRem_reconstruct (x y : Int) (nonzero : y ≠ 0) :
    (intQuoRem x y).map (fun (q, r) => y * q + r) = some x := by
  simp [intQuoRem, nonzero, Int.mul_tdiv_add_tmod]

theorem intDiv_negative_dividend :
    intDiv (-7) 3 = some (-3) ∧ intDiv (-7) (-3) = some 3 := ⟨rfl, rfl⟩

theorem intMod_negative_dividend :
    intMod (-7) 3 = some 2 ∧ intMod (-7) (-3) = some 2 := ⟨rfl, rfl⟩

theorem intDivMod_reconstruct (x y : Int) (nonzero : y ≠ 0) :
    (intDivMod x y).map (fun (q, r) => y * q + r) = some x := by
  have rebuild := Int.mul_tdiv_add_tmod x y
  by_cases less : x.tmod y < 0
  · by_cases negative : y < 0
    · simp [intDivMod, intQuoRem, nonzero, less, negative, Int.mul_add]
      omega
    · simp [intDivMod, intQuoRem, nonzero, less, negative, Int.mul_sub]
      omega
  · simp [intDivMod, intQuoRem, nonzero, less, rebuild]

theorem intMod_euclidean_bounds (x y r : Int) (defined : intMod x y = some r) :
    0 ≤ r ∧ r < (if y < 0 then -y else y) := by
  have nonzero : y ≠ 0 := by
    intro h
    simp [intMod, intDivMod, intQuoRem, h] at defined
  by_cases negative : y < 0
  · have pos : 0 < -y := by omega
    have lower := Int.lt_tmod_of_pos x pos
    have upper := Int.tmod_lt_of_pos x pos
    rw [Int.tmod_neg] at lower upper
    by_cases remNeg : x.tmod y < 0
    · simp [intMod, intDivMod, intQuoRem, nonzero, negative, remNeg] at defined ⊢
      omega
    · simp [intMod, intDivMod, intQuoRem, nonzero, negative, remNeg] at defined ⊢
      omega
  · have pos : 0 < y := by omega
    have lower := Int.lt_tmod_of_pos x pos
    have upper := Int.tmod_lt_of_pos x pos
    by_cases remNeg : x.tmod y < 0
    · simp [intMod, intDivMod, intQuoRem, nonzero, negative, remNeg] at defined ⊢
      omega
    · simp [intMod, intDivMod, intQuoRem, nonzero, negative, remNeg] at defined ⊢
      omega

end GoLean.Math.Big
