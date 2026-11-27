import Init.Data.BitVec.Lemmas

namespace GoLean.Math

/-- Go 1.27.0 reinterprets the binary64 representation as a signed 64-bit
integer and compares it with zero. The input is the representation returned by
Float64bits.

https://github.com/golang/go/blob/go1.27.0/src/math/signbit.go#L7-L10
https://github.com/golang/go/blob/go1.27.0/src/math/unsafe.go#L32-L35 -/
def signbit (repr : BitVec 64) : Bool :=
  decide (repr.toInt < 0)

theorem signbit_spec (repr : BitVec 64) : signbit repr = repr.getLsbD 63 := by
  change decide (repr.toInt < 0) = repr.getLsbD (64 - 1)
  rw [← BitVec.msb_eq_getLsbD_last]
  cases h : repr.msb with
  | false => simp [BitVec.toInt_nonneg_of_msb_false h]
  | true => simp [BitVec.toInt_neg_of_msb_true h]

end GoLean.Math
