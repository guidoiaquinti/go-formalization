import Init.Data.BitVec.Lemmas
import Init.Data.Int.DivMod.Lemmas
import Init.Data.Nat.Sqrt.Lemmas

namespace GoLean.Math

/-- One digit per moving bit, returning its residual and twice-root
    accumulator along with the root. -/
@[irreducible] private def sqrtDigits : Nat → Nat → Nat → Nat → Nat × Nat × Nat
  | 0, rest, root, twice => (rest, root, twice)
  | steps + 1, rest, root, twice =>
    let bit := 2 ^ steps
    let trial := twice + bit
    if trial ≤ rest then
      sqrtDigits steps (2 * (rest - trial)) (root + bit) (trial + bit)
    else sqrtDigits steps (2 * rest) root twice

private def sqrtInput (exponent fraction : Nat) : Nat × Int :=
  if exponent = 0 then
    let shift := 52 - fraction.log2
    (fraction <<< shift, 1 - (shift : Int) - 1023)
  else (2 ^ 52 + fraction, (exponent : Int) - 1023)

private def sqrtRadicand (significand : Nat) (exponent : Int) : Nat :=
  if exponent % 2 = 1 then 2 * significand else significand

private def assembleSqrt (remainder root : Nat) (exponent : Int) : Nat :=
  let rounded := if remainder = 0 then root else root + root % 2
  rounded / 2 + (exponent / 2 + 1022).toNat * 2 ^ 52

@[irreducible] private def sqrtCore (steps radicand : Nat) (unbiased : Int) : Nat :=
  let (remainder, root, _) := sqrtDigits (steps + 1) (2 * radicand) 0 0
  assembleSqrt remainder root unbiased

/-- Models the generic moving-bit loop. The leading-set-bit calculation
    replaces the equivalent subnormal normalization loop. Natural
    intermediates represent the nonnegative word shifts and updates.
    None represents a NaN result by class.
    https://github.com/golang/go/blob/go1.27.0/src/math/sqrt.go#L85-L145 -/
def sqrtWord (word : BitVec 64) : Option Nat :=
  let bits := word.toNat
  let magnitude := bits % 2 ^ 63
  let exponent := magnitude / 2 ^ 52
  let fraction := magnitude % 2 ^ 52
  if magnitude = 0 ∨ bits = 2047 * 2 ^ 52 then some bits
  else if bits ≥ 2 ^ 63 ∨ exponent = 2047 then none
  else
    let (significand, unbiased) := sqrtInput exponent fraction
    let radicand := sqrtRadicand significand unbiased
    some (sqrtCore 53 radicand unbiased)

private theorem squareStep (q bit : Nat) :
    (q + bit) * (q + bit) = q * q + (2 * q + bit) * bit := by
  calc
    _ = q * q + (q * bit + q * bit + bit * bit) := by
      simp only [Nat.mul_add, Nat.mul_comm]
      ac_rfl
    _ = _ := by rw [Nat.add_mul, Nat.mul_assoc, Nat.two_mul]

private theorem sqrtDigits_floor (steps rest root target : Nat)
    (balance : root * root + rest * 2 ^ steps = target)
    (ceiling : target < (root + 2 ^ (steps + 1)) * (root + 2 ^ (steps + 1))) :
    let result := sqrtDigits (steps + 1) rest root (2 * root)
    result.2.1 * result.2.1 + result.1 / 2 = target ∧
      result.1 % 2 = 0 ∧ target < (result.2.1 + 1) * (result.2.1 + 1) := by
  induction steps generalizing rest root target with
  | zero =>
    simp only [Nat.pow_zero, Nat.mul_one] at balance
    change target < (root + 1 + 1) * (root + 1 + 1) at ceiling
    by_cases accepts : 2 * root + 1 ≤ rest
    · simp only [sqrtDigits, Nat.pow_zero, ite_eq_left accepts]
      rw [squareStep]
      omega
    · simp only [sqrtDigits, Nat.pow_zero, ite_eq_right accepts]
      rw [squareStep]
      omega
  | succ steps ih =>
    let bit := 2 ^ (steps + 1)
    let trial := 2 * root + bit
    have nextScale : bit = 2 * 2 ^ steps := by
      change 2 ^ (steps + 1) = _
      rw [Nat.pow_succ, Nat.mul_comm]
    have upperScale : 2 ^ (steps + 1 + 1) = bit + bit := by
      rw [Nat.pow_succ, Nat.mul_two]
    rw [upperScale] at ceiling
    change root * root + rest * bit = target at balance
    rw [sqrtDigits]
    change (let result := if trial ≤ rest then
        sqrtDigits (steps + 1) (2 * (rest - trial)) (root + bit) (trial + bit)
      else sqrtDigits (steps + 1) (2 * rest) root (2 * root)
      result.2.1 * result.2.1 + result.1 / 2 = target ∧
        result.1 % 2 = 0 ∧ target < (result.2.1 + 1) * (result.2.1 + 1))
    by_cases accepts : trial ≤ rest
    · rw [ite_eq_left accepts, show trial + bit = 2 * (root + bit) by dsimp [trial]; omega]
      apply ih
      · have scale : 2 * (rest - trial) * 2 ^ steps = (rest - trial) * bit := by
          rw [nextScale]
          ac_rfl
        rw [scale, squareStep]
        change root * root + trial * bit + (rest - trial) * bit = target
        rw [Nat.add_assoc, ← Nat.add_mul, Nat.add_sub_of_le accepts, balance]
      · simpa only [Nat.add_assoc] using ceiling
    · rw [ite_eq_right accepts]
      apply ih
      · have scale : 2 * rest * 2 ^ steps = rest * bit := by
          rw [nextScale]
          ac_rfl
        rw [scale, balance]
      · change target < (root + bit) * (root + bit)
        rw [squareStep, ← balance]
        exact Nat.add_lt_add_left
          (Nat.mul_lt_mul_of_pos_right (Nat.lt_of_not_ge accepts) (Nat.two_pow_pos _)) _

private theorem sqrtDigits_normal (steps radicand : Nat) (bound : radicand < 2 ^ (steps + 1)) :
    let target := radicand * 2 ^ (steps + 1)
    let digits := sqrtDigits (steps + 1) (2 * radicand) 0 0
    digits.2.1 = target.sqrt ∧ digits.1 = 2 * (target - target.sqrt * target.sqrt) := by
  let target := radicand * 2 ^ (steps + 1)
  let digits := sqrtDigits (steps + 1) (2 * radicand) 0 0
  let root := digits.2.1
  have balance : 0 * 0 + (2 * radicand) * 2 ^ steps = target := by
    unfold target
    rw [Nat.zero_mul, Nat.zero_add, Nat.pow_succ]
    ac_rfl
  have ceiling : target < (0 + 2 ^ (steps + 1)) * (0 + 2 ^ (steps + 1)) := by
    simpa only [Nat.zero_add] using
      Nat.mul_lt_mul_of_pos_right bound (Nat.two_pow_pos (steps + 1))
  have proof := sqrtDigits_floor steps (2 * radicand) 0 target balance ceiling
  change root * root + digits.1 / 2 = target ∧ digits.1 % 2 = 0 ∧
    target < (root + 1) * (root + 1) at proof
  have rootEq : root = target.sqrt := by
    have sqLo := Nat.sqrt_le target
    have sqHi := Nat.lt_succ_sqrt target
    change target < (target.sqrt + 1) * (target.sqrt + 1) at sqHi
    by_cases tooSmall : root < target.sqrt
    · have between := Nat.mul_self_le_mul_self (show root + 1 ≤ target.sqrt by omega)
      omega
    · by_cases tooLarge : target.sqrt < root
      · have between := Nat.mul_self_le_mul_self (show target.sqrt + 1 ≤ root by omega)
        omega
      · omega
  refine ⟨rootEq, ?_⟩
  have balance := proof.1
  have even := proof.2.1
  rw [← rootEq]
  change digits.1 = 2 * (target - root * root)
  omega

private theorem sqrtInput_bound (e f : Nat) (nonzero : e ≠ 0 ∨ f ≠ 0)
    (fractionBound : f < 2 ^ 52) : (sqrtInput e f).1 < 2 ^ 53 := by
  unfold sqrtInput
  by_cases subnormal : e = 0
  · rw [ite_eq_left subnormal]
    have fNonzero : f ≠ 0 := nonzero.resolve_left (not_not_intro subnormal)
    have topBound : f.log2 < 52 := (Nat.log2_lt fNonzero).mpr fractionBound
    have shifted := Nat.mul_lt_mul_of_pos_right (Nat.lt_log2_self (n := f))
      (Nat.two_pow_pos (52 - f.log2))
    rw [← Nat.pow_add, show f.log2 + 1 + (52 - f.log2) = 53 by omega] at shifted
    simpa only [Nat.shiftLeft_eq] using shifted
  · rw [ite_eq_right subnormal]
    dsimp only
    omega

private theorem sqrtCore_spec (steps radicand : Nat) (unbiased : Int)
    (bound : radicand < 2 ^ (steps + 1)) :
    sqrtCore steps radicand unbiased =
      let scaled := radicand * 2 ^ (steps + 1)
      let root := scaled.sqrt
      root / 2 + (if root % 2 = 1 ∧ root * root < scaled then 1 else 0) +
        (unbiased / 2 + 1022).toNat * 2 ^ 52 := by
  unfold sqrtCore
  dsimp only
  let digits := sqrtDigits (steps + 1) (2 * radicand) 0 0
  let scaled := radicand * 2 ^ (steps + 1)
  let root := scaled.sqrt
  change assembleSqrt digits.1 digits.2.1 unbiased =
    root / 2 + (if root % 2 = 1 ∧ root * root < scaled then 1 else 0) +
      (unbiased / 2 + 1022).toNat * 2 ^ 52
  have spec := sqrtDigits_normal steps radicand bound
  change digits.2.1 = root ∧ digits.1 = 2 * (scaled - root * root) at spec
  rw [spec.1, spec.2, assembleSqrt]
  have lower : root * root ≤ scaled := Nat.sqrt_le scaled
  by_cases exact : root * root = scaled
  · simp [exact]
  · rw [ite_eq_right (by omega : ¬ 2 * (scaled - root * root) = 0)]
    by_cases odd : root % 2 = 1
    · rw [ite_eq_left (show root % 2 = 1 ∧ root * root < scaled from ⟨odd, by omega⟩)]
      omega
    · rw [ite_eq_right (fun h => odd h.1)]
      omega

/-- The moving-bit loop computes the integer square root of the scaled
    radicand. Its discarded low digit rounds upward only when that digit
    is set and the integer square leaves a nonzero residue. -/
theorem sqrtWord_spec (word : BitVec 64) :
    sqrtWord word =
      let bits := word.toNat
      let magnitude := bits % 2 ^ 63
      let e := magnitude / 2 ^ 52
      let f := magnitude % 2 ^ 52
      if magnitude = 0 ∨ bits = 2047 * 2 ^ 52 then some bits
      else if bits ≥ 2 ^ 63 ∨ e = 2047 then none
      else
        let normalized := sqrtInput e f
        let scaled := sqrtRadicand normalized.1 normalized.2 * 2 ^ 54
        let root := scaled.sqrt
        let addUnit := if root % 2 = 1 ∧ root * root < scaled then 1 else 0
        some (root / 2 + addUnit + (normalized.2 / 2 + 1022).toNat * 2 ^ 52) := by
  unfold sqrtWord
  dsimp only
  by_cases zeroOrInf : word.toNat % 2 ^ 63 = 0 ∨ word.toNat = 2047 * 2 ^ 52
  · rw [ite_eq_left zeroOrInf, ite_eq_left zeroOrInf]
  · rw [ite_eq_right zeroOrInf, ite_eq_right zeroOrInf]
    by_cases negativeOrNaN : 2 ^ 63 ≤ word.toNat ∨ word.toNat % 2 ^ 63 / 2 ^ 52 = 2047
    · rw [ite_eq_left negativeOrNaN, ite_eq_left negativeOrNaN]
    · rw [ite_eq_right negativeOrNaN, ite_eq_right negativeOrNaN]
      let magnitude := word.toNat % 2 ^ 63
      let e := magnitude / 2 ^ 52
      let f := magnitude % 2 ^ 52
      let normalized := sqrtInput e f
      let rad := sqrtRadicand normalized.1 normalized.2
      change some (sqrtCore 53 rad normalized.2) = _
      have nonzero : e ≠ 0 ∨ f ≠ 0 := by
        have mNonzero : magnitude ≠ 0 := fun h => zeroOrInf (.inl h)
        have parts := Nat.mod_add_div' magnitude (2 ^ 52)
        omega
      have normalBound := sqrtInput_bound e f nonzero (Nat.mod_lt _ (by decide))
      have radBound : rad < 2 ^ 54 := by
        unfold rad sqrtRadicand
        change (if normalized.2 % 2 = 1 then 2 * normalized.1 else normalized.1) < 2 ^ 54
        change normalized.1 < 2 ^ 53 at normalBound
        split <;> omega
      rw [sqrtCore_spec 53 rad normalized.2 radBound]

end GoLean.Math
