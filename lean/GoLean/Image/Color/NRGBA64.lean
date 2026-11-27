import Init.Data.BitVec.Lemmas

namespace GoLean.Image.Color

private def premultiply16 (channel alpha : BitVec 16) : BitVec 32 :=
  (channel.zeroExtend 32 * alpha.zeroExtend 32) / 0xffff

/-- Go widens the input words to uint32, then multiplies and divides each RGB
    channel. The alpha is widened unchanged.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/color.go#L78-L96 -/
def nrgba64RGBA (r g b a : BitVec 16) : BitVec 32 × BitVec 32 × BitVec 32 × BitVec 32 :=
  (premultiply16 r a, premultiply16 g a, premultiply16 b a, a.zeroExtend 32)

private theorem widen16 (x : BitVec 16) : (x.zeroExtend 32).toNat = x.toNat :=
  BitVec.toNat_setWidth_of_le (by decide)

private theorem premultiply16_spec (channel alpha : BitVec 16) :
    (premultiply16 channel alpha).toNat = channel.toNat * alpha.toNat / 65535 := by
  have bound : channel.toNat * alpha.toNat < 2 ^ 32 := BitVec.toNat_mul_toNat_lt
  unfold premultiply16
  rw [BitVec.toNat_udiv, BitVec.toNat_mul, widen16, widen16, Nat.mod_eq_of_lt bound]
  rfl

theorem nrgba64RGBA_spec (r g b a : BitVec 16) :
    (let (red, green, blue, alpha) := nrgba64RGBA r g b a
     (red.toNat, green.toNat, blue.toNat, alpha.toNat)) =
      (r.toNat * a.toNat / 65535, g.toNat * a.toNat / 65535,
       b.toNat * a.toNat / 65535, a.toNat) := by
  simp only [nrgba64RGBA, premultiply16_spec, widen16]

end GoLean.Image.Color
