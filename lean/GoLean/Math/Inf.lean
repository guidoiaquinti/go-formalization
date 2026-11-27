import Init.Data.BitVec.Basic

namespace GoLean.Math

/-- Go selects a word based on the signed comparison and reinterprets that word
    as a float64. The model returns the resulting binary64 representation.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L7-L28
    https://github.com/golang/go/blob/go1.27.0/src/math/unsafe.go#L37-L41 -/
def inf (sign : Int) : BitVec 64 :=
  if sign ≥ 0 then 0x7ff0000000000000 else 0xfff0000000000000

theorem inf_spec (sign : Int) :
    (inf sign).toNat = (if sign < 0 then 2 ^ 63 else 0) + 2047 * 2 ^ 52 := by
  by_cases h : sign < 0
  · rw [inf, ite_eq_right (by omega : ¬ sign ≥ 0), ite_eq_left h]
    decide
  · rw [inf, ite_eq_left (by omega : sign ≥ 0), ite_eq_right h]
    decide

end GoLean.Math
