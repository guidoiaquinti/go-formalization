import Init.Data.BitVec.Lemmas

namespace GoLean.Image.Color

private def replicateByte (channel : BitVec 8) : BitVec 32 :=
  let wide := channel.zeroExtend 32
  wide ||| (wide <<< 8)

/-- Go widens and repeats the byte in each of the four channels.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/color.go#L21-L40 -/
def rgbaRGBA (r g b a : BitVec 8) : BitVec 32 × BitVec 32 × BitVec 32 × BitVec 32 :=
  (replicateByte r, replicateByte g, replicateByte b, replicateByte a)

private theorem replicateByte_spec : ∀ x : BitVec 8,
    (replicateByte x).toNat = x.toNat * 257 := by
  decide

theorem rgbaRGBA_spec (r g b a : BitVec 8) :
    (let (red, green, blue, alpha) := rgbaRGBA r g b a
     (red.toNat, green.toNat, blue.toNat, alpha.toNat)) =
    (r.toNat * 257, g.toNat * 257, b.toNat * 257, a.toNat * 257) := by
  simp only [rgbaRGBA, replicateByte_spec]

end GoLean.Image.Color
