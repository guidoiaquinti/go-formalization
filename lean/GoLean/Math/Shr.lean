import Init.Data.BitVec.Lemmas

namespace GoLean.Math

/-- Go shifts the high word down and ORs together three contributions to
    the low word. When a subtraction on the native unsigned shift count
    wraps, the resulting count is at least 64, and that term is zero.
    The model substitutes a zero term in each wrapping case.
    https://github.com/golang/go/blob/go1.27.0/src/math/fma.go#L33-L37 -/
def shr (high low : BitVec 64) (n : Nat) : BitVec 64 × BitVec 64 :=
  let fromHighLeft := if n ≤ 64 then high <<< (64 - n) else 0
  let fromHighRight := if n ≥ 64 then high >>> (n - 64) else 0
  (high >>> n, (low >>> n) ||| fromHighLeft ||| fromHighRight)

private theorem bitOutside (word : BitVec w) (h : w ≤ i) :
    word.getLsbD i = false := by
  cases bit : word.getLsbD i
  · rfl
  · exact False.elim (Nat.not_lt_of_ge h (BitVec.lt_of_getLsbD bit))

/-- The two returned limbs form the unsigned right shift of the combined
    128-bit value, including zero when the count reaches 128. -/
theorem shr_spec (high low : BitVec 64) (n : Nat) :
    (shr high low n).1 ++ (shr high low n).2 = (high ++ low) >>> n := by
  apply BitVec.eq_of_getLsbD_eq_iff.mpr
  intro i _
  simp only [shr, BitVec.getLsbD_ushiftRight, BitVec.getLsbD_append]
  by_cases inLow : i < 64
  · simp only [ite_eq_left inLow, BitVec.getLsbD_or, BitVec.getLsbD_ushiftRight]
    by_cases staysLow : n + i < 64
    · have small : n < 64 := by omega
      simp [staysLow, Nat.le_of_lt small, show ¬ 64 ≤ n by omega,
        inLow, show i < 64 - n by omega]
    · rw [ite_eq_right staysLow, bitOutside low (by omega)]
      by_cases small : n < 64
      · simp [Nat.le_of_lt small, show ¬ 64 ≤ n by omega, inLow,
          -BitVec.getLsbD_eq_getElem, show ¬ i < 64 - n by omega, show i - (64 - n) = n + i - 64 by omega]
      · by_cases middle : n = 64
        · subst n
          simp
        · simp [show ¬ n ≤ 64 by omega, show 64 ≤ n by omega,
            show (n - 64) + i = n + i - 64 by omega]
  · simp [inLow, show ¬ n + i < 64 by omega, show n + (i - 64) = n + i - 64 by omega]

end GoLean.Math
