import Init.Data.BitVec.Basic

namespace GoLean.Math

/-- A binary64 value represented by its stored encoding. -/
structure Binary64 where
  encoding : BitVec 64

/-- Read the same word at the address of the floating-point input.
    https://github.com/golang/go/blob/go1.27.0/src/math/unsafe.go#L32-L35 -/
def float64Bits (f : Binary64) : BitVec 64 := f.encoding

/-- Interpret the same stored word as a floating-point value.
    https://github.com/golang/go/blob/go1.27.0/src/math/unsafe.go#L37-L41 -/
def float64FromBits (bits : BitVec 64) : Binary64 := ⟨bits⟩

theorem float64Bits_spec (f : Binary64) :
    float64Bits f = f.encoding := rfl

theorem float64FromBits_spec (bits : BitVec 64) :
    (float64FromBits bits).encoding = bits := rfl

theorem float64Bits_roundtrip (f : Binary64) :
    float64FromBits (float64Bits f) = f := by
  cases f
  rfl

theorem float64FromBits_roundtrip (bits : BitVec 64) :
    float64Bits (float64FromBits bits) = bits := rfl

end GoLean.Math
