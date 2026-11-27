import Init.Data.BitVec.Basic

namespace GoLean.Math

/-- A binary32 value represented by its stored encoding. -/
structure Binary32 where
  encoding : BitVec 32

/-- Read the same word at the address of the floating-point input.
    https://github.com/golang/go/blob/go1.27.0/src/math/unsafe.go#L21-L24 -/
def float32Bits (f : Binary32) : BitVec 32 := f.encoding

/-- Interpret the same stored word as a floating-point value.
    https://github.com/golang/go/blob/go1.27.0/src/math/unsafe.go#L26-L30 -/
def float32FromBits (bits : BitVec 32) : Binary32 := ⟨bits⟩

theorem float32Bits_spec (f : Binary32) :
    float32Bits f = f.encoding := rfl

theorem float32FromBits_spec (bits : BitVec 32) :
    (float32FromBits bits).encoding = bits := rfl

theorem float32Bits_roundtrip (f : Binary32) :
    float32FromBits (float32Bits f) = f := by
  cases f
  rfl

theorem float32FromBits_roundtrip (bits : BitVec 32) :
    float32Bits (float32FromBits bits) = bits := rfl

end GoLean.Math
