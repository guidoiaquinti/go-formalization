import Init.Data.BitVec.Lemmas

namespace GoLean.Math

/-- Go masks the magnitude and sign words separately and reinterprets their OR.
    https://github.com/golang/go/blob/go1.27.0/src/math/copysign.go#L7-L11
    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L15
    https://github.com/golang/go/blob/go1.27.0/src/math/unsafe.go#L32-L41 -/
def copysign (magnitude sign : BitVec 64) : BitVec 64 :=
  let signMask : BitVec 64 := 0x8000000000000000
  (magnitude &&& ~~~signMask) ||| (sign &&& signMask)

/-- Keep the low 63 bits of the first word and the top bit of the second. -/
theorem copysign_spec (magnitude sign : BitVec 64) :
    (copysign magnitude sign).toNat =
      magnitude.toNat % (2 ^ 63) + (sign.toNat / (2 ^ 63)) * (2 ^ 63) := by
  rw [copysign, BitVec.toNat_or, BitVec.toNat_and, BitVec.toNat_and]
  change (magnitude.toNat &&& (2 ^ 63 - 1)) ||| (sign.toNat &&& 2 ^ 63) = _
  rw [Nat.and_two_pow_sub_one_eq_mod]
  have bound := sign.isLt
  have small : sign.toNat / 2 ^ 63 = 0 ∨ sign.toNat / 2 ^ 63 = 1 := by omega
  have quotient : (sign.toNat &&& 2 ^ 63) / 2 ^ 63 = sign.toNat / 2 ^ 63 := by
    rw [Nat.and_div_two_pow]
    rcases small with h | h <;> rw [h] <;> decide
  have remainder : (sign.toNat &&& 2 ^ 63) % 2 ^ 63 = 0 := by
    rw [Nat.and_mod_two_pow]
    simp
  have whole : sign.toNat &&& 2 ^ 63 = sign.toNat / (2 ^ 63) * 2 ^ 63 := by omega
  rw [whole]
  rcases small with h | h
  · simp [h]
  · simp only [h, Nat.one_mul]
    exact Nat.or_two_pow_eq_add_of_lt (Nat.mod_lt _ (by decide))

end GoLean.Math
