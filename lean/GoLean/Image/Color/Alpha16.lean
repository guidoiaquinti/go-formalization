import Init.Data.BitVec.Lemmas

namespace GoLean.Image.Color

/-- Go widens the alpha field and uses the resulting alpha for all four channels.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/color.go#L109-L117 -/
def alpha16RGBA (input : BitVec 16) : BitVec 32 × BitVec 32 × BitVec 32 × BitVec 32 :=
  let alpha := input.zeroExtend 32
  (alpha, alpha, alpha, alpha)

theorem alpha16RGBA_spec (input : BitVec 16) :
    (let (r, g, b, a) := alpha16RGBA input
     (r.toNat, g.toNat, b.toNat, a.toNat)) =
      (input.toNat, input.toNat, input.toNat, input.toNat) := by
  have bound : input.toNat < 2 ^ 16 := input.isLt
  simp [alpha16RGBA]
  omega

end GoLean.Image.Color
