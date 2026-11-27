import Init.Data.BitVec.Lemmas

namespace GoLean.Math.RandV2

/-- Go 1.27.0 math/rand/v2.absInt32. Signed negation and the unsigned
conversion retain the same 32 bits, including for the minimum signed input. -/
def absInt32 (x : BitVec 32) : BitVec 32 :=
  if x.msb then -x else x

theorem absInt32_sign (x : BitVec 32) :
    (if x.msb then -(absInt32 x) else absInt32 x) = x := by
  cases h : x.msb <;> simp [absInt32, h, BitVec.neg_neg]

theorem absInt32_min : absInt32 (0x80000000 : BitVec 32) = 0x80000000 := by decide

end GoLean.Math.RandV2
