import Init.Data.BitVec.Lemmas

namespace GoLean.Image.Color

/-- Go 1.27.0 widens the grayscale field to uint32 and returns it for all three
    color channels with constant opaque alpha.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/color.go#L130-L138 -/
def gray16RGBA (y : BitVec 16) : BitVec 32 × BitVec 32 × BitVec 32 × BitVec 32 :=
  let wide := y.zeroExtend 32
  (wide, wide, wide, 0xffff)

/-- Each RGB channel has exactly the input's unsigned value, and alpha is opaque. -/
theorem gray16RGBA_spec (y : BitVec 16) :
    (let (r, g, b, a) := gray16RGBA y
     (r.toNat, g.toNat, b.toNat, a.toNat)) = (y.toNat, y.toNat, y.toNat, 65535) := by
  have bound : y.toNat < 2 ^ 16 := y.isLt
  simp [gray16RGBA]
  omega

end GoLean.Image.Color
