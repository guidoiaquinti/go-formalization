import Init.Data.BitVec.Lemmas

namespace GoLean.Image.Color

/-- Go widens the gray byte and ORs it with its left shift by eight.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/color.go#L119-L128 -/
def grayRGBA (y : BitVec 8) : BitVec 32 × BitVec 32 × BitVec 32 × BitVec 32 :=
  let wide := y.zeroExtend 32
  let channel := wide ||| (wide <<< 8)
  (channel, channel, channel, 0xffff)

theorem grayRGBA_spec : ∀ y : BitVec 8,
    (let (r, g, b, a) := grayRGBA y
     (r.toNat, g.toNat, b.toNat, a.toNat)) =
      (y.toNat * 257, y.toNat * 257, y.toNat * 257, 65535) := by
  decide

end GoLean.Image.Color
