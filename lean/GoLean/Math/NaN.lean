import Init.Data.BitVec.Basic

namespace GoLean.Math

/-- Go reinterprets uvnan as float64. The model returns that binary64 word.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L7-L8
    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L30-L31
    https://github.com/golang/go/blob/go1.27.0/src/math/unsafe.go#L37-L41 -/
def nan : BitVec 64 := 0x7ff8000000000001

/-- A positive sign, an all-ones exponent, the quiet bit, and payload one. -/
theorem nan_spec : nan.toNat = 2047 * 2 ^ 52 + 2 ^ 51 + 1 := by decide


end GoLean.Math
