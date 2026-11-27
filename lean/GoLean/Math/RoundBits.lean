import Init.Data.BitVec.Lemmas

namespace GoLean.Math.RoundBits

theorem clearLow (bits k : Nat) (hbits : bits < 2 ^ 64) (hk : k ≤ 64) :
    bits &&& (2 ^ 64 - 2 ^ k) = bits / 2 ^ k * 2 ^ k := by
  have power : 2 ^ (64 - k) * 2 ^ k = (2 : Nat) ^ 64 := by
    rw [← Nat.pow_add, Nat.sub_add_cancel hk]
  have mask : 2 ^ 64 - 2 ^ k = (2 ^ (64 - k) - 1) * 2 ^ k := by
    rw [Nat.sub_mul, Nat.one_mul, power]
  have highBound : bits / 2 ^ k < 2 ^ (64 - k) := by
    rw [Nat.div_lt_iff_lt_mul (Nat.two_pow_pos k), power]
    exact hbits
  have upper : (bits &&& (2 ^ 64 - 2 ^ k)) / 2 ^ k = bits / 2 ^ k := by
    rw [Nat.and_div_two_pow, mask, Nat.mul_div_cancel _ (Nat.two_pow_pos k),
      Nat.and_two_pow_sub_one_of_lt_two_pow highBound]
  have lower : (bits &&& (2 ^ 64 - 2 ^ k)) % 2 ^ k = 0 := by
    rw [Nat.and_mod_two_pow, mask, Nat.mul_mod_left, Nat.and_zero]
  have combine := Nat.mod_add_div' (bits &&& (2 ^ 64 - 2 ^ k)) (2 ^ k)
  rw [upper, lower, Nat.zero_add] at combine
  exact combine.symm

theorem shiftedPow (n k : Nat) (hk : k ≤ n) :
    2 ^ n >>> k = 2 ^ (n - k) := by
  rw [Nat.shiftRight_eq_div_pow, ← Nat.pow_sub_mul_pow 2 hk,
    Nat.mul_div_cancel _ (Nat.two_pow_pos k)]

theorem shiftedMask (n k : Nat) (hk : k ≤ n) :
    (2 ^ n - 1) >>> k = 2 ^ (n - k) - 1 := by
  have product := Nat.pow_sub_mul_pow 2 hk
  have lower : (2 ^ (n - k) - 1) * 2 ^ k = 2 ^ n - 2 ^ k := by
    rw [Nat.sub_mul, Nat.one_mul, product]
  have power := Nat.two_pow_pos k
  rw [Nat.shiftRight_eq_div_pow, Nat.div_eq_iff power, lower]
  have positive := Nat.two_pow_pos (n - k)
  have bound : 2 ^ k ≤ 2 ^ n := by
    calc
      2 ^ k = 1 * 2 ^ k := by simp
      _ ≤ 2 ^ (n - k) * 2 ^ k := Nat.mul_le_mul_right _ positive
      _ = 2 ^ n := product
  omega

end GoLean.Math.RoundBits
