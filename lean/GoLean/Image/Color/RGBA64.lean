import Init.Data.BitVec.Lemmas

namespace GoLean.Image.Color

/-- Go widens each unsigned 16-bit field to a 32-bit channel.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/color.go#L42-L53 -/
def rgba64RGBA (r g b a : BitVec 16) : BitVec 32 × BitVec 32 × BitVec 32 × BitVec 32 :=
  (r.zeroExtend 32, g.zeroExtend 32, b.zeroExtend 32, a.zeroExtend 32)

theorem rgba64RGBA_spec (r g b a : BitVec 16) :
    (let (red, green, blue, alpha) := rgba64RGBA r g b a
     (red.toNat, green.toNat, blue.toNat, alpha.toNat)) =
    (r.toNat, g.toNat, b.toNat, a.toNat) := by
  simp only [rgba64RGBA, BitVec.zeroExtend, BitVec.toNat_setWidth_of_le (by decide : 16 ≤ 32)]

end GoLean.Image.Color
