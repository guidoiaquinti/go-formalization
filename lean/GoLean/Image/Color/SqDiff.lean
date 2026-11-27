import Init.Data.BitVec.Lemmas

namespace GoLean.Image.Color

/-- The Go subtraction and multiplication both wrap at 32 bits.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/color.go#L310-L339 -/
def sqDiff (x y : BitVec 32) : BitVec 32 :=
  let d := x - y
  (d * d) >>> 2

private theorem sqDiff_ordered (x y : BitVec 32) (hxy : y ≤ x) (hx : x.toNat ≤ 65535) :
    (sqDiff x y).toNat = (x.toNat - y.toNat) * (x.toNat - y.toNat) / 4 := by
  have hd : x.toNat - y.toNat ≤ 65535 := Nat.le_trans (Nat.sub_le _ _) hx
  have squareBound : (x.toNat - y.toNat) * (x.toNat - y.toNat) < 2 ^ 32 :=
    Nat.lt_of_le_of_lt (Nat.mul_le_mul hd hd) (by decide)
  rw [sqDiff, BitVec.toNat_ushiftRight, BitVec.toNat_mul_of_lt, BitVec.toNat_sub_of_le hxy]
  · rw [Nat.shiftRight_eq_div_pow]
  · rw [BitVec.toNat_sub_of_le hxy]
    exact squareBound

theorem sqDiff_spec (x y : BitVec 32) (hx : x.toNat ≤ 65535) (hy : y.toNat ≤ 65535) :
    (sqDiff x y).toNat = (x.toNat - y.toNat + (y.toNat - x.toNat)) ^ 2 / 4 := by
  by_cases hxy : y ≤ x
  · have hnat : y.toNat ≤ x.toNat := hxy
    rw [sqDiff_ordered x y hxy hx, Nat.sub_eq_zero_of_le hnat, Nat.add_zero, Nat.pow_two]
  · have hnat : x.toNat ≤ y.toNat := by
      change ¬ y.toNat ≤ x.toNat at hxy
      omega
    have swap : x - y = -(y - x) := by rw [BitVec.neg_sub, BitVec.sub_eq_add_neg, BitVec.add_comm]
    have symmetric : sqDiff x y = sqDiff y x := by
      simp only [sqDiff]
      rw [swap, BitVec.neg_mul_neg]
    rw [symmetric, sqDiff_ordered y x hnat hy, Nat.sub_eq_zero_of_le hnat, Nat.zero_add, Nat.pow_two]

end GoLean.Image.Color
