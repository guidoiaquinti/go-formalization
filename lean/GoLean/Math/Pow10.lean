import Init.Data.Float.Float

namespace GoLean.Math

private def pow10Small : Array Float := #[
  1.0e0, 1.0e1, 1.0e2, 1.0e3, 1.0e4, 1.0e5, 1.0e6, 1.0e7, 1.0e8, 1.0e9,
  1.0e10, 1.0e11, 1.0e12, 1.0e13, 1.0e14, 1.0e15, 1.0e16, 1.0e17, 1.0e18, 1.0e19,
  1.0e20, 1.0e21, 1.0e22, 1.0e23, 1.0e24, 1.0e25, 1.0e26, 1.0e27, 1.0e28, 1.0e29,
  1.0e30, 1.0e31]

private def pow10Positive : Array Float := #[
  1.0e0, 1.0e32, 1.0e64, 1.0e96, 1.0e128, 1.0e160, 1.0e192, 1.0e224,
  1.0e256, 1.0e288]

private def pow10Negative : Array Float := #[
  1.0e0, 1.0e-32, 1.0e-64, 1.0e-96, 1.0e-128, 1.0e-160, 1.0e-192, 1.0e-224,
  1.0e-256, 1.0e-288, 1.0e-320]

/-- Go divides the exponent magnitude into a block of 32 and a remainder,
    then multiplies or divides the rounded binary64 table values.
    https://github.com/golang/go/blob/go1.27.0/src/math/pow10.go#L7-L48 -/
def pow10 (n : Int) : Float :=
  if positive : 0 ≤ n ∧ n ≤ 308 then
    pow10Positive[n.toNat / 32]'(by change n.toNat / 32 < 10; omega) *
      pow10Small[n.toNat % 32]'(by change n.toNat % 32 < 32; omega)
  else if negative : -323 ≤ n ∧ n < 0 then
    pow10Negative[(-n).toNat / 32]'(by change (-n).toNat / 32 < 11; omega) /
      pow10Small[(-n).toNat % 32]'(by change (-n).toNat % 32 < 32; omega)
  else if n > 0 then Float.ofBits 0x7ff0000000000000 else 0

theorem pow10_tails (n : Int) :
    (n < -323 → pow10 n = 0) ∧
      (308 < n → pow10 n = Float.ofBits 0x7ff0000000000000) := by
  constructor
  · intro tiny
    rw [pow10, dite_eq_right (by omega), dite_eq_right (by omega), ite_eq_right (by omega)]
  · intro huge
    rw [pow10, dite_eq_right (by omega), dite_eq_right (by omega), ite_eq_left (by omega)]

theorem pow10_positive_block (block offset : Nat) (small : offset < 32)
    (supported : block * 32 + offset ≤ 308) :
    pow10 ((↑(block * 32 + offset) : Int)) =
      pow10Positive[block]'(by change block < 10; omega) *
        pow10Small[offset]'(by exact small) := by
  rw [pow10, dite_eq_left (by omega)]
  simp only [Int.toNat_natCast, Nat.mul_comm block, Nat.mul_add_mod_self_left,
    Nat.mod_eq_of_lt small, Nat.mul_add_div (by decide : 0 < (32 : Nat)),
    Nat.div_eq_of_lt small, Nat.add_zero]
  rfl

theorem pow10_negative_block (block offset : Nat) (small : offset < 32)
    (nonzero : 0 < block * 32 + offset) (supported : block * 32 + offset ≤ 323) :
    pow10 (-(↑(block * 32 + offset) : Int)) =
      pow10Negative[block]'(by change block < 11; omega) /
        pow10Small[offset]'(by exact small) := by
  rw [pow10, dite_eq_right (by omega), dite_eq_left (by omega)]
  simp only [Int.neg_neg, Int.toNat_natCast, Nat.mul_comm block, Nat.mul_add_mod_self_left,
    Nat.mod_eq_of_lt small, Nat.mul_add_div (by decide : 0 < (32 : Nat)),
    Nat.div_eq_of_lt small, Nat.add_zero]
  rfl

end GoLean.Math
