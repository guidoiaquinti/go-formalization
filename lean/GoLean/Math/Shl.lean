import Init.Data.BitVec.Lemmas

namespace GoLean.Math

/-- Go shifts the low word up and ORs together three contributions to the
    high word. When unsigned subtraction wraps, the resulting shift count
    is at least 64 and the term is zero. Counts of at least 128 make every
    term zero and can be handled before performing any shift.
    https://github.com/golang/go/blob/go1.27.0/src/math/fma.go#L27-L31 -/
def shl (high low : BitVec 64) (n : Nat) : BitVec 64 × BitVec 64 :=
  if n ≥ 128 then (0, 0)
  else
    let fromLowRight := if n ≤ 64 then low >>> (64 - n) else 0
    let fromLowLeft := if n ≥ 64 then low <<< (n - 64) else 0
    ((high <<< n) ||| fromLowRight ||| fromLowLeft, low <<< n)

private theorem bitOutside (word : BitVec w) (h : w ≤ i) :
    word.getLsbD i = false := by
  cases bit : word.getLsbD i
  · rfl
  · exact False.elim (Nat.not_lt_of_ge h (BitVec.lt_of_getLsbD bit))

/-- Rejoining the two returned limbs gives the fixed-width 128-bit left
    shift. Bits above the 128-bit destination are discarded. -/
theorem shl_spec (high low : BitVec 64) (n : Nat) :
    (shl high low n).1 ++ (shl high low n).2 = (high ++ low) <<< n := by
  unfold shl
  by_cases large : 128 ≤ n
  · rw [ite_eq_left large, BitVec.shiftLeft_eq_zero large]
    rfl
  · rw [ite_eq_right large]
    dsimp only
    apply BitVec.eq_of_getLsbD_eq_iff.mpr
    intro i bound
    simp only [BitVec.getLsbD_append, BitVec.getLsbD_shiftLeft, bound, decide_true,
      Bool.true_and]
    by_cases inLow : i < 64
    · simp [inLow, show i - n < 64 by omega]
    · rw [ite_eq_right inLow]
      have offsetBound : i - 64 < 64 := by omega
      by_cases small : n < 64
      · simp only [ite_eq_left (Nat.le_of_lt small), ite_eq_right (show ¬ n ≥ 64 by omega),
          BitVec.getLsbD_or, BitVec.getLsbD_shiftLeft, BitVec.getLsbD_ushiftRight,
          offsetBound, decide_true, Bool.true_and,
          show 64 - n + (i - 64) = i - n by omega]
        by_cases inHigh : 64 ≤ i - n
        · rw [bitOutside low inHigh]
          simp [show ¬ i - n < 64 by omega, show ¬ i < n by omega,
            show ¬ i - 64 < n by omega, show i - 64 - n = i - n - 64 by omega,
            -BitVec.getLsbD_eq_getElem]
        · simp [show i - n < 64 by omega, show ¬ i < n by omega,
            show i - 64 < n by omega]
      · by_cases middle : n = 64
        · subst n
          simp [offsetBound, inLow, -BitVec.getLsbD_eq_getElem]
        · simp only [ite_eq_right (show ¬ n ≤ 64 by omega),
            ite_eq_left (show n ≥ 64 by omega),
            BitVec.getLsbD_or, BitVec.getLsbD_shiftLeft,
            offsetBound, decide_true, Bool.true_and,
            show i - 64 < n by omega, decide_true, Bool.not_true, Bool.false_and,
            Bool.false_or, show (i - 64 < n - 64) ↔ i < n by omega,
            show i - 64 - (n - 64) = i - n by omega]
          simp [show i - n < 64 by omega]

end GoLean.Math
