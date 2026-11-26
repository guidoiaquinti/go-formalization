import Init.Data.BitVec.Lemmas

namespace GoLean.Math.Bits

inductive UintWidth where
  | bits32
  | bits64

abbrev UintWidth.bits : UintWidth → Nat
  | .bits32 => 32
  | .bits64 => 64

inductive RemOutcome (w : Nat) where
  | divideByZero
  | value (r : BitVec w)

/-- Go 1.27.0 math/bits.Rem, for any three uints of the platform's width.
    The 32-bit branch takes the remainder of the widened dividend. The 64-bit
    branch first reduces hi modulo y, so even hi ≥ y cannot overflow Div64's
    quotient. Both branches panic on a zero divisor.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L573
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L580-L597 -/
def remUint (width : UintWidth) (hi lo y : BitVec width.bits) : RemOutcome width.bits :=
  if y.toNat = 0 then .divideByZero
  else
    let high := match width with
      | .bits32 => hi.toNat
      | .bits64 => hi.toNat % y.toNat
    .value (BitVec.ofNat width.bits ((high * 2 ^ width.bits + lo.toNat) % y.toNat))

/-- Correct one estimated quotient digit as in the Go Div64 source.
    An exhausted estimate is reported to the differential harness. -/
private def correctDigit (yn1 yn0 next : BitVec 64) :
    Nat → BitVec 64 → BitVec 64 → Option (BitVec 64)
  | fuel, q, r =>
    let base : BitVec 64 := 0x100000000
    if q < base ∧ q * yn0 ≤ base * r + next then some q
    else
      match fuel with
      | 0 => none
      | passes + 1 =>
        let q := q - 1
        let r := r + yn1
        if base ≤ r then some q else correctDigit yn1 yn0 next passes q r

/-- For a normalized two-digit divisor, the estimate from its high
    digit is at most two above the true next quotient digit. This
    arithmetic bound concerns the estimate before machine-word correction. -/
theorem div64_quotient_digit_bound (high low partialDividend next : Nat)
    (normalized : 2^31 ≤ high) (lowFits : low < 2^32)
    (partialDividendFits : partialDividend < high * 2^32 + low) (nextFits : next < 2^32) :
    let estimate := partialDividend / high
    let digit := (partialDividend * 2^32 + next) / (high * 2^32 + low)
    digit ≤ estimate ∧ estimate ≤ digit + 2 := by
  let base := 2^32
  let divisor := high * base + low
  let dividend := partialDividend * base + next
  let estimate := partialDividend / high
  let digit := dividend / divisor
  change digit ≤ estimate ∧ estimate ≤ digit + 2
  have highPos : 0 < high := by omega
  have divisorPos : 0 < divisor := by dsimp [divisor, base]; omega
  have partialDividendBound : partialDividend < (estimate + 1) * high :=
    (Nat.div_lt_iff_lt_mul highPos).1 (Nat.lt_succ_self estimate)
  have estimateBound : estimate < base + 2 := by
    apply (Nat.div_lt_iff_lt_mul highPos).2
    rw [Nat.add_mul, Nat.mul_comm base high]
    dsimp [base]
    omega
  constructor
  · have digitBound : dividend < (estimate + 1) * divisor := calc
      dividend < ((estimate + 1) * high) * base := by
        dsimp [dividend, base] at *
        omega
      _ ≤ (estimate + 1) * divisor := by
        rw [Nat.mul_assoc]
        exact Nat.mul_le_mul_left _ (Nat.le_add_right ..)
    have := (Nat.div_lt_iff_lt_mul divisorPos).2 digitBound
    omega
  · by_cases small : estimate < 2
    · exact Nat.le_trans (Nat.le_of_lt small) (Nat.le_add_left ..)
    let m := estimate - 2
    have mFits : m ≤ base := by dsimp [m]; omega
    have exactEstimate : m + 2 = estimate := by dsimp [m]; omega
    have lowBound : m * low ≤ (2 * high) * base :=
      Nat.le_trans (Nat.mul_le_mul mFits (Nat.le_of_lt lowFits))
        (Nat.mul_le_mul_right base (by dsimp [base]; omega))
    have fits : m * divisor ≤ dividend := calc
      m * divisor = (m * high) * base + m * low := by
        simp [divisor, Nat.mul_add, Nat.mul_assoc]
      _ ≤ ((m + 2) * high) * base := by
        rw [Nat.add_mul, Nat.add_mul]
        exact Nat.add_le_add_left lowBound _
      _ = (estimate * high) * base := by rw [exactEstimate]
      _ ≤ partialDividend * base := Nat.mul_le_mul_right base (Nat.div_mul_le_self ..)
      _ ≤ dividend := Nat.le_add_right ..
    have := (Nat.le_div_iff_mul_le divisorPos).2 fits
    dsimp [m] at this
    omega

/-- The correction guard accepts exactly a one-limb estimate whose
    product does not exceed the current dividend. The range hypotheses
    ensure that neither side of the word comparison has wrapped. -/
theorem div64_quotient_digit_guard (q r high low next : BitVec 64)
    (qFits : q.toNat ≤ 2^32 + 1) (rFits : r.toNat < 2^32)
    (lowFits : low.toNat < 2^32) (nextFits : next.toNat < 2^32) :
    (q < 0x100000000 ∧ q * low ≤ 0x100000000 * r + next) ↔
      q.toNat < 2^32 ∧
        q.toNat * (high.toNat * 2^32 + low.toNat) ≤
          (q.toNat * high.toNat + r.toNat) * 2^32 + next.toNat := by
  have productFits : q.toNat * low.toNat < 2^64 := by
    have bound := Nat.mul_le_mul qFits (show low.toNat ≤ 2^32 - 1 by omega)
    omega
  have shiftedFits : 2^32 * r.toNat < 2^64 := by omega
  have sumFits : 2^32 * r.toNat + next.toNat < 2^64 := by omega
  have base : (0x100000000 : BitVec 64).toNat = 2^32 := by decide
  simp only [BitVec.lt_def, BitVec.le_def, BitVec.toNat_mul, BitVec.toNat_add, base,
    Nat.mod_eq_of_lt productFits, Nat.mod_eq_of_lt shiftedFits, Nat.mod_eq_of_lt sumFits,
    Nat.mul_add, Nat.add_mul, ← Nat.mul_assoc]
  rw [Nat.mul_comm (2^32) r.toNat]
  omega

/-- Starting with an estimate at most `fuel` too high, the word-level
    correction loop returns the exact digit. -/
private theorem correctDigit_eq_of_invariant
    (high low next : BitVec 64) (partialDividend : Nat)
    (highPos : 0 < high.toNat) (highFits : high.toNat < 2^32)
    (lowFits : low.toNat < 2^32) (nextFits : next.toNat < 2^32)
    (partialFits : partialDividend < high.toNat * 2^32 + low.toNat)
    (fuel : Nat) (q r : BitVec 64)
    (decompose : q.toNat * high.toNat + r.toNat = partialDividend)
    (qFits : q.toNat ≤ 2^32 + 1) (rFits : r.toNat < 2^32)
    (lower : (partialDividend * 2^32 + next.toNat) /
      (high.toNat * 2^32 + low.toNat) ≤ q.toNat)
    (upper : q.toNat ≤ (partialDividend * 2^32 + next.toNat) /
      (high.toNat * 2^32 + low.toNat) + fuel) :
    correctDigit high low next fuel q r =
      some (BitVec.ofNat 64 ((partialDividend * 2^32 + next.toNat) /
        (high.toNat * 2^32 + low.toNat))) := by
  let divisor := high.toNat * 2^32 + low.toNat
  let dividend := partialDividend * 2^32 + next.toNat
  let digit := dividend / divisor
  change digit ≤ q.toNat at lower
  change q.toNat ≤ digit + fuel at upper
  change correctDigit high low next fuel q r = some (BitVec.ofNat 64 digit)
  have pos : 0 < divisor := by dsimp [divisor]; omega
  have digitFits : digit < 2^32 := by
    apply (Nat.div_lt_iff_lt_mul pos).2
    rw [Nat.mul_comm (2^32) divisor]
    dsimp [dividend, divisor]
    omega
  induction fuel generalizing q r with
  | zero =>
    have eq : q.toNat = digit := by omega
    have guard := (div64_quotient_digit_guard q r high low next qFits rFits lowFits nextFits).2
      ⟨by omega, by rw [decompose, eq]; exact Nat.div_mul_le_self ..⟩
    have word : q = BitVec.ofNat 64 digit := by rw [← eq]; simp
    rw [correctDigit, ite_eq_left guard, word]
  | succ fuel ih =>
    by_cases guard : q < 0x100000000 ∧ q * low ≤ 0x100000000 * r + next
    · have accepted := (div64_quotient_digit_guard q r high low next qFits rFits lowFits nextFits).1 guard
      rw [decompose] at accepted
      have le := (Nat.le_div_iff_mul_le pos).2 accepted.2
      have eq : q.toNat = digit := Nat.le_antisymm le lower
      have word : q = BitVec.ofNat 64 digit := by rw [← eq]; simp
      rw [correctDigit, ite_eq_left guard, word]
    · have greater : digit < q.toNat := by
        apply Nat.lt_of_not_ge
        intro eq
        have same : q.toNat = digit := Nat.le_antisymm eq lower
        apply guard
        apply (div64_quotient_digit_guard q r high low next qFits rFits lowFits nextFits).2
        refine ⟨by omega, ?_⟩
        rw [decompose, same]
        exact Nat.div_mul_le_self ..
      have qPos : 0 < q.toNat := Nat.lt_of_le_of_lt (Nat.zero_le digit) greater
      let q' := q - 1
      let r' := r + high
      have qNat : q'.toNat = q.toNat - 1 := by
        apply BitVec.toNat_sub_of_not_usubOverflow
        simp [BitVec.usubOverflow]
        omega
      have rNat : r'.toNat = r.toNat + high.toNat := by
        apply BitVec.toNat_add_of_not_uaddOverflow
        simp [BitVec.uaddOverflow]
        omega
      have nextDecompose : q'.toNat * high.toNat + r'.toNat = partialDividend := by
        rw [rNat, ← Nat.add_assoc, Nat.add_right_comm (q'.toNat * _) r.toNat,
          ← Nat.succ_mul, show q'.toNat.succ = q.toNat by change q'.toNat + 1 = q.toNat; omega]
        exact decompose
      have nextLower : digit ≤ q'.toNat := by omega
      by_cases stop : (0x100000000 : BitVec 64) ≤ r'
      · have stopNat : 2^32 ≤ r'.toNat := by
          simpa [BitVec.le_def] using stop
        have qSmall : q'.toNat < 2^32 := by
          apply Nat.lt_of_not_ge
          intro large
          have le := Nat.mul_le_mul_right high.toNat large
          rw [Nat.mul_comm (2^32) high.toNat] at le
          omega
        have productBound : q'.toNat * low.toNat ≤ r'.toNat * 2^32 :=
          Nat.mul_le_mul (by omega) (by omega)
        have le : q'.toNat ≤ digit := by
          apply (Nat.le_div_iff_mul_le pos).2
          dsimp [divisor, dividend]
          rw [← nextDecompose, Nat.mul_add, Nat.add_mul, ← Nat.mul_assoc]
          omega
        have word : q' = BitVec.ofNat 64 digit := by
          rw [← Nat.le_antisymm le nextLower]
          simp
        simp only [correctDigit, ite_eq_right guard, show 0x100000000 ≤ r + high from stop,
          ite_true]
        exact congrArg some word
      · have rSmall : r'.toNat < 2^32 := by
          simpa [BitVec.le_def] using stop
        simp only [correctDigit, ite_eq_right guard,
          ite_eq_right (show ¬ 0x100000000 ≤ r + high from stop)]
        exact ih q' r' nextDecompose (by omega) rSmall nextLower (by omega)

/-- The two correction steps suffice when division starts from the
    source's estimated quotient and remainder. -/
theorem div64_correctDigit (partialDividend high low next : BitVec 64)
    (normalized : 2^31 ≤ high.toNat) (highFits : high.toNat < 2^32)
    (lowFits : low.toNat < 2^32) (nextFits : next.toNat < 2^32)
    (partialFits : partialDividend.toNat < high.toNat * 2^32 + low.toNat) :
    let q := partialDividend / high
    correctDigit high low next 2 q (partialDividend - q * high) =
      some (BitVec.ofNat 64 ((partialDividend.toNat * 2^32 + next.toNat) /
        (high.toNat * 2^32 + low.toNat))) := by
  let q := partialDividend / high
  let r := partialDividend - q * high
  change correctDigit high low next 2 q r = _
  have highPos : 0 < high.toNat := by omega
  have quotient : q.toNat = partialDividend.toNat / high.toNat := BitVec.toNat_udiv
  have productBound : q.toNat * high.toNat ≤ partialDividend.toNat := by
    rw [quotient]
    exact Nat.div_mul_le_self ..
  have product : (q * high).toNat = q.toNat * high.toNat := by
    rw [BitVec.toNat_mul, Nat.mod_eq_of_lt (Nat.lt_of_le_of_lt productBound partialDividend.isLt)]
  have remainder : r.toNat = partialDividend.toNat % high.toNat := by
    rw [BitVec.toNat_sub_of_not_usubOverflow (by simpa [BitVec.usubOverflow, product] using productBound),
      product, quotient]
    have := Nat.div_add_mod partialDividend.toNat high.toNat
    rw [Nat.mul_comm high.toNat] at this
    omega
  have fits : q.toNat ≤ 2^32 + 1 := by
    have bound : partialDividend.toNat < (2^32 + 2) * high.toNat := by
      rw [Nat.add_mul, Nat.mul_comm (2^32) high.toNat]
      omega
    have := (Nat.div_lt_iff_lt_mul highPos).2 bound
    omega
  have estimate := div64_quotient_digit_bound high.toNat low.toNat partialDividend.toNat next.toNat
    normalized lowFits partialFits nextFits
  have digits := Nat.div_add_mod partialDividend.toNat high.toNat
  rw [Nat.mul_comm high.toNat] at digits
  apply correctDigit_eq_of_invariant high low next partialDividend.toNat highPos highFits
    lowFits nextFits partialFits 2 q r
  · rw [quotient, remainder, digits]
  · exact fits
  · exact remainder ▸ Nat.lt_trans (Nat.mod_lt _ highPos) highFits
  · exact quotient ▸ estimate.1
  · exact quotient ▸ estimate.2

/-- A corrected quotient digit leaves the exact remainder even if the
    intermediate products wrap in unsigned arithmetic. -/
private theorem div64_digit_remainder (partialDividend next divisor : BitVec 64) :
    let n := partialDividend.toNat * 2^32 + next.toNat
    let q := BitVec.ofNat 64 (n / divisor.toNat)
    partialDividend * 0x100000000 + next - q * divisor =
      BitVec.ofNat 64 (n % divisor.toNat) := by
  let n := partialDividend.toNat * 2^32 + next.toNat
  let p := n / divisor.toNat * divisor.toNat
  have decompose : n % divisor.toNat + p = n := Nat.mod_add_div' ..
  have newDividend : partialDividend * 0x100000000 + next = BitVec.ofNat 64 n := by
    dsimp [n]
    rw [BitVec.ofNat_add, BitVec.ofNat_mul]
    simp
  have quotientProduct :
      BitVec.ofNat 64 (n / divisor.toNat) * divisor = BitVec.ofNat 64 p := by
    dsimp [p]
    rw [BitVec.ofNat_mul]
    simp
  change partialDividend * 0x100000000 + next -
    BitVec.ofNat 64 (n / divisor.toNat) * divisor = _
  rw [newDividend, quotientProduct, ← decompose, BitVec.ofNat_add, BitVec.add_sub_cancel]

private def div64NormalizedRemainder (un32 un10 y : BitVec 64) : Option (BitVec 64) := do
  let base : BitVec 64 := 0x100000000
  let mask := base - 1
  let yn1 := y >>> (32 : Nat)
  let yn0 := y &&& mask
  let un1 := un10 >>> (32 : Nat)
  let un0 := un10 &&& mask
  let q1 := un32 / yn1
  let q1 ← correctDigit yn1 yn0 un1 2 q1 (un32 - q1 * yn1)
  let un21 := un32 * base + un1 - q1 * y
  let q0 := un21 / yn1
  let q0 ← correctDigit yn1 yn0 un0 2 q0 (un21 - q0 * yn1)
  return un21 * base + un0 - q0 * y

private theorem div64NormalizedRemainder_eq (un32 un10 y : BitVec 64)
    (normalized : 2^63 ≤ y.toNat) (hiFits : un32.toNat < y.toNat) :
    div64NormalizedRemainder un32 un10 y =
      some (BitVec.ofNat 64 ((un32.toNat * 2^64 + un10.toNat) % y.toNat)) := by
  let mask : BitVec 64 := 0xffffffff
  let yn1 := y >>> (32 : Nat)
  let yn0 := y &&& mask
  let un1 := un10 >>> (32 : Nat)
  let un0 := un10 &&& mask
  have highDigit (x : BitVec 64) : (x >>> (32 : Nat)).toNat = x.toNat / 2^32 := by
    simp [Nat.shiftRight_eq_div_pow]
  have lowDigit (x : BitVec 64) : (x &&& mask).toNat = x.toNat % 2^32 := by
    have maskNat : mask.toNat = 2^32 - 1 := by decide
    rw [BitVec.toNat_and, maskNat, Nat.and_two_pow_sub_one_eq_mod]
  have combine (x : BitVec 64) :
      (x >>> (32 : Nat)).toNat * 2^32 + (x &&& mask).toNat = x.toNat := by
    rw [highDigit, lowDigit, Nat.add_comm, Nat.mod_add_div']
  have yDigits : yn1.toNat * 2^32 + yn0.toNat = y.toNat := combine y
  have upperFits (x : BitVec 64) : (x >>> (32 : Nat)).toNat < 2^32 := by
    rw [highDigit]
    have := x.isLt
    omega
  have lowerFits (x : BitVec 64) : (x &&& mask).toNat < 2^32 := by
    rw [lowDigit]
    exact Nat.mod_lt _ (by decide)
  have highNormalized : 2^31 ≤ yn1.toNat := by
    change 2^31 ≤ (y >>> (32 : Nat)).toNat
    rw [highDigit]
    omega
  have pos : 0 < y.toNat := by omega
  have correct (p next : BitVec 64) (np : next.toNat < 2^32) (pf : p.toNat < y.toNat) :
      correctDigit yn1 yn0 next 2 (p / yn1) (p - (p / yn1) * yn1) =
        some (BitVec.ofNat 64 ((p.toNat * 2^32 + next.toNat) / y.toNat)) := by
    have h := div64_correctDigit p yn1 yn0 next highNormalized (upperFits y)
      (lowerFits y) np (by rwa [yDigits])
    rwa [yDigits] at h
  let first := un32.toNat * 2^32 + un1.toNat
  let middle := BitVec.ofNat 64 (first % y.toNat)
  have middleNat : middle.toNat = first % y.toNat := by
    rw [BitVec.toNat_ofNat, Nat.mod_eq_of_lt (Nat.lt_trans (Nat.mod_lt _ pos) y.isLt)]
  have midFits : middle.toNat < y.toNat := by rw [middleNat]; exact Nat.mod_lt _ pos
  have fused : (middle.toNat * 2^32 + un0.toNat) % y.toNat =
      (un32.toNat * 2^64 + un10.toNat) % y.toNat := by
    rw [middleNat, Nat.add_mod_eq_add_mod_right _ (Nat.mod_mul_mod ..)]
    dsimp [first]
    rw [Nat.add_mul, Nat.mul_assoc, Nat.add_assoc, combine un10]
  change (correctDigit yn1 yn0 un1 2 (un32 / yn1) (un32 - (un32 / yn1) * yn1)).bind
    (fun q1 =>
      let m := un32 * 0x100000000 + un1 - q1 * y
      (correctDigit yn1 yn0 un0 2 (m / yn1) (m - (m / yn1) * yn1)).bind
        (fun q0 => some (m * 0x100000000 + un0 - q0 * y))) = _
  rw [correct un32 un1 (upperFits un10) hiFits]
  rw [Option.bind_some]
  rw [div64_digit_remainder un32 un1 y]
  change (correctDigit yn1 yn0 un0 2 (middle / yn1) (middle - (middle / yn1) * yn1)).bind
    (fun q0 => some (middle * 0x100000000 + un0 - q0 * y)) = _
  rw [correct middle un0 (lowerFits un10) midFits]
  rw [Option.bind_some]
  rw [div64_digit_remainder middle un0 y, fused]

/-- Remainder from the pinned Div64 source path, after Rem64 reduces the
    high word. The caller supplies a positive divisor and hi < y.
    The quotient is discarded by Rem64. -/
private def div64SourceRemainder (hi lo y : BitVec 64) : Option (BitVec 64) := do
  if hi = 0 then return lo % y
  let shift := 64 - (Nat.log2 y.toNat + 1)
  let un32 := (hi <<< shift) ||| (lo >>> (64 - shift))
  let un10 := lo <<< shift
  let r ← div64NormalizedRemainder un32 un10 (y <<< shift)
  return r >>> shift

private theorem div64SourceRemainder_eq (hi lo y : BitVec 64)
    (yPos : 0 < y.toNat) (highFits : hi.toNat < y.toNat) :
    div64SourceRemainder hi lo y =
      some (BitVec.ofNat 64 ((hi.toNat * 2^64 + lo.toNat) % y.toNat)) := by
  by_cases zero : hi = 0
  · subst hi
    simp [div64SourceRemainder, BitVec.umod_def]
  let shift := 64 - (y.toNat.log2 + 1)
  let scale := 2 ^ shift
  let remaining := 2 ^ (64 - shift)
  let un32 := (hi <<< shift) ||| (lo >>> (64 - shift))
  let un10 := lo <<< shift
  let scaledY := y <<< shift
  have logBound := (Nat.log2_eq_iff (Nat.ne_of_gt yPos)).1 (rfl : y.toNat.log2 = _)
  have logFits := (Nat.log2_lt (Nat.ne_of_gt yPos)).2 y.isLt
  have scalePos : 0 < scale := Nat.two_pow_pos _
  have splitPower : remaining * scale = 2^64 := by
    rw [← Nat.pow_add]
    congr 1
    dsimp [shift]
    omega
  have logScale : 2 ^ y.toNat.log2 * scale = 2^63 := by
    rw [← Nat.pow_add]
    congr 1
    dsimp [shift]
    omega
  have divisorFits : y.toNat * scale < 2^64 := by
    calc
      _ < 2 ^ (y.toNat.log2 + 1) * scale := Nat.mul_lt_mul_of_pos_right logBound.2 scalePos
      _ = _ := by rw [← Nat.pow_add]; congr 1; dsimp [shift]; omega
  have highScaled : hi.toNat * scale < 2^64 :=
    Nat.lt_trans (Nat.mul_lt_mul_of_pos_right highFits scalePos) divisorFits
  have yNat : scaledY.toNat = y.toNat * scale := by
    rw [BitVec.toNat_shiftLeft, Nat.shiftLeft_eq, Nat.mod_eq_of_lt divisorFits]
  have loNat : un10.toNat = lo.toNat * scale % 2^64 := by
    rw [BitVec.toNat_shiftLeft, Nat.shiftLeft_eq]
  have spillFits : lo.toNat / remaining < scale := by
    apply (Nat.div_lt_iff_lt_mul (Nat.two_pow_pos _)).2
    rw [Nat.mul_comm, splitPower]
    exact lo.isLt
  have hiNat : un32.toNat = hi.toNat * scale + lo.toNat / remaining := by
    rw [BitVec.toNat_or, BitVec.toNat_shiftLeft, Nat.mod_eq_of_lt
      (show hi.toNat <<< shift < 2^64 by rw [Nat.shiftLeft_eq]; exact highScaled),
      BitVec.toNat_ushiftRight, Nat.shiftRight_eq_div_pow,
      ← Nat.shiftLeft_add_eq_or_of_lt spillFits, Nat.shiftLeft_eq]
  have hiNorm : un32.toNat < scaledY.toNat := by
    rw [hiNat, yNat]
    have le := Nat.mul_le_mul_right scale (Nat.succ_le_of_lt highFits)
    rw [Nat.succ_mul] at le
    omega
  have yNorm : 2^63 ≤ scaledY.toNat := by
    rw [yNat, ← logScale]
    exact Nat.mul_le_mul_right _ logBound.1
  have spill : lo.toNat * scale / 2^64 = lo.toNat / remaining := by
    rw [← splitPower, Nat.mul_div_mul_right _ _ scalePos]
  have loAssembled : lo.toNat / remaining * 2^64 + un10.toNat = lo.toNat * scale := by
    rw [loNat, ← spill, Nat.add_comm, Nat.mod_add_div']
  have dividend : un32.toNat * 2^64 + un10.toNat =
      (hi.toNat * 2^64 + lo.toNat) * scale := by
    rw [hiNat, Nat.add_mul, Nat.add_assoc, loAssembled, Nat.add_mul, Nat.mul_right_comm hi.toNat]
  rw [div64SourceRemainder, ite_eq_right zero]
  change (div64NormalizedRemainder un32 un10 scaledY).bind (fun r => some (r >>> shift)) = _
  rw [div64NormalizedRemainder_eq un32 un10 scaledY yNorm hiNorm, dividend, yNat,
    Nat.mul_mod_mul_right]
  rw [Option.bind_some, Option.some.injEq]
  apply BitVec.eq_of_toNat_eq
  have remainderFits : (hi.toNat * 2^64 + lo.toNat) % y.toNat * scale < 2^64 :=
    Nat.lt_trans (Nat.mul_lt_mul_of_pos_right (Nat.mod_lt _ yPos) scalePos) divisorFits
  rw [BitVec.toNat_ushiftRight, Nat.shiftRight_eq_div_pow, BitVec.toNat_ofNat,
    Nat.mod_eq_of_lt remainderFits, Nat.mul_div_cancel _ scalePos]
  rw [BitVec.toNat_ofNat, Nat.mod_eq_of_lt (Nat.lt_trans (Nat.mod_lt _ yPos) y.isLt)]

/-- Source-path model used alongside the proved arithmetic abstraction. -/
def remUintSource (width : UintWidth) (hi lo y : BitVec width.bits) :
    Option (RemOutcome width.bits) :=
  if y.toNat = 0 then some .divideByZero
  else
    match width with
    | .bits32 =>
      let n := (hi.zeroExtend 64 <<< 32) ||| lo.zeroExtend 64
      some (.value ((n % y.zeroExtend 64).truncate 32))
    | .bits64 => (div64SourceRemainder (hi % y) lo y).map .value

/-- The source and arithmetic models both report division by zero,
    before any width-specific arithmetic. -/
theorem remUintSource_eq_zero (width : UintWidth) (hi lo : BitVec width.bits) :
    remUintSource width hi lo 0 = some (remUint width hi lo 0) := by
  simp [remUintSource, remUint]

/-- The widened 32-bit source operations implement the arithmetic model
    on the full input domain. -/
theorem remUintSource_eq_bits32 (hi lo y : BitVec 32) :
    remUintSource .bits32 hi lo y = some (remUint .bits32 hi lo y) := by
  have widen (x : BitVec 32) : (x.zeroExtend 64).toNat = x.toNat :=
    BitVec.toNat_setWidth_of_le (by decide)
  have shiftFits : hi.toNat <<< 32 < 2^64 := by
    rw [Nat.shiftLeft_eq]
    have := hi.isLt
    omega
  have assembled :
      (((hi.zeroExtend 64) <<< (32 : Nat)) ||| lo.zeroExtend 64).toNat =
        hi.toNat * 2^32 + lo.toNat := by
    rw [BitVec.toNat_or, BitVec.toNat_shiftLeft, widen, widen, Nat.mod_eq_of_lt shiftFits,
      ← Nat.shiftLeft_add_eq_or_of_lt lo.isLt, Nat.shiftLeft_eq]
  by_cases hz : y.toNat = 0
  · have eq : y = 0 := BitVec.eq_of_toNat_eq (by simpa using hz)
    subst y
    exact remUintSource_eq_zero .bits32 hi lo
  · simp only [remUintSource, remUint, UintWidth.bits, ite_eq_right hz,
      Option.some.injEq, RemOutcome.value.injEq]
    apply BitVec.eq_of_toNat_eq
    change (_ % y.zeroExtend 64).toNat % 2^32 = _
    rw [BitVec.toNat_umod, widen, assembled, BitVec.toNat_ofNat]

/-- The source-shaped remainder model agrees with the arithmetic
    model at both word widths, including when the divisor is zero. -/
theorem remUintSource_eq (width : UintWidth) (hi lo y : BitVec width.bits) :
    remUintSource width hi lo y = some (remUint width hi lo y) := by
  cases width with
  | bits32 => exact remUintSource_eq_bits32 hi lo y
  | bits64 =>
    by_cases hz : y.toNat = 0
    · have eq : y = 0 := BitVec.eq_of_toNat_eq (by simpa using hz)
      subst y
      exact remUintSource_eq_zero .bits64 hi lo
    · have pos : 0 < y.toNat := Nat.pos_of_ne_zero hz
      have highFits : (hi % y).toNat < y.toNat := Nat.mod_lt _ pos
      simp only [remUintSource, remUint, UintWidth.bits, ite_eq_right hz]
      rw [div64SourceRemainder_eq (hi % y) lo y pos highFits]
      rfl

/-- The panic requires a zero divisor. Every return is the unique Euclidean
    remainder of the original two-word dividend, including when hi ≥ y. -/
theorem remUint_spec (width : UintWidth) (hi lo y : BitVec width.bits) :
    match remUint width hi lo y with
    | .divideByZero => y.toNat = 0
    | .value r => r < y ∧ ∃ q : Nat, hi.toNat * 2 ^ width.bits + lo.toNat = y.toNat * q + r.toNat := by
  by_cases hz : y.toNat = 0
  · simp [remUint, hz]
  · have pos : 0 < y.toNat := by omega
    let n := hi.toNat * 2 ^ width.bits + lo.toNat
    have congruent : (((hi.toNat % y.toNat) * 2 ^ width.bits + lo.toNat) % y.toNat) =
        n % y.toNat := Nat.add_mod_eq_add_mod_right _ (Nat.mod_mul_mod ..)
    have hr := Nat.mod_lt n pos
    have fits := Nat.lt_trans hr y.isLt
    dsimp [n] at fits congruent
    cases width <;>
      simp only [remUint, hz, ↓reduceIte, BitVec.lt_def, BitVec.toNat_ofNat, congruent,
        Nat.mod_eq_of_lt fits] <;>
      exact ⟨hr, n / y.toNat, (Nat.div_add_mod ..).symm⟩

end GoLean.Math.Bits
