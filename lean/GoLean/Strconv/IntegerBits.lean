import Init.Data.BitVec.Lemmas

namespace GoLean.Strconv

/-- The pinned internal/strconv lower helper unconditionally sets bit five. -/
def lower (x : BitVec 8) : BitVec 8 := x ||| 32

/-- The helper accepts zero and tests the bit pattern of a native signed int. -/
def isPowerOfTwo (x : BitVec width) : Bool := (x &&& (x - 1)) == 0

theorem lower_fixedBits (x : BitVec 8) :
    lower x &&& 223 = x &&& 223 ∧ lower x &&& 32 = 32 := by
  unfold lower
  constructor
  · rw [BitVec.and_or_distrib_right]
    exact BitVec.or_zero
  · apply BitVec.eq_of_getLsbD_eq_iff.mpr
    intro i _
    simp only [BitVec.getLsbD_and, BitVec.getLsbD_or]
    cases (32 : BitVec 8).getLsbD i <;> simp

theorem powerOfTwo_zero (width : Nat) : isPowerOfTwo (0 : BitVec width) = true := by
  simp [isPowerOfTwo]

theorem powerOfTwo_min32 : isPowerOfTwo (0x80000000 : BitVec 32) = true := by decide

theorem powerOfTwo_min64 : isPowerOfTwo (0x8000000000000000 : BitVec 64) = true := by decide

end GoLean.Strconv
