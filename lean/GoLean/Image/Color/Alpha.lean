import Init.Data.BitVec.Lemmas

namespace GoLean.Image.Color

/-- Go widens the alpha field and uses the resulting alpha for all four channels.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/color.go#L98-L107 -/
def alphaRGBA (input : BitVec 8) : BitVec 32 × BitVec 32 × BitVec 32 × BitVec 32 :=
  let wide := input.zeroExtend 32
  let alpha := wide ||| (wide <<< 8)
  (alpha, alpha, alpha, alpha)

theorem alphaRGBA_spec : ∀ input : BitVec 8,
    (let (r, g, b, a) := alphaRGBA input
     (r.toNat, g.toNat, b.toNat, a.toNat)) =
      (input.toNat * 257, input.toNat * 257, input.toNat * 257, input.toNat * 257) := by
  decide

end GoLean.Image.Color
