import Init.Data.BitVec.Lemmas

namespace GoLean.Math.Bits

/-- Go 1.27.0, src/math/bits/bits.go, lines 276–278.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L276-L278 -/
def reverseBytes16 (x : BitVec 16) : BitVec 16 :=
  (x >>> 8) ||| (x <<< 8)

/-- The result consists of the original low byte followed by the original high byte. -/
theorem reverseBytes16_spec (x : BitVec 16) :
    reverseBytes16 x = x.extractLsb' 0 8 ++ x.extractLsb' 8 8 := by
  apply BitVec.eq_of_getLsbD_eq
  intro i hi
  rw [BitVec.getLsbD_append (x := x.extractLsb' 0 8) (y := x.extractLsb' 8 8)]
  by_cases low : i < 8
  · simp [reverseBytes16, low]
  · have outside := BitVec.getLsbD_of_ge x (8 + i) (by omega)
    have offset : i - 8 < 8 := by omega
    simp [reverseBytes16, low, hi, outside, offset]
    rfl

theorem reverseBytes16_involutive (x : BitVec 16) :
    reverseBytes16 (reverseBytes16 x) = x := by
  rw [reverseBytes16_spec, reverseBytes16_spec, BitVec.extractLsb'_append_eq_right,
    BitVec.extractLsb'_append_eq_left, BitVec.extractLsb'_append_extractLsb']

end GoLean.Math.Bits
