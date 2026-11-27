import GoLean.Image.Color.CMYK

namespace GoLean.Image.Color

/-- Go computes the same widened channel arithmetic as CMYK.RGBA and then
    narrows each channel after shifting right by eight.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/ycbcr.go#L335-L360 -/
def cmykToRGB (c m y k : BitVec 8) : BitVec 8 × BitVec 8 × BitVec 8 :=
  let (r, g, b, _) := cmykRGBA c m y k
  ((r >>> 8).truncate 8, (g >>> 8).truncate 8, (b >>> 8).truncate 8)

private theorem narrow_spec (word : BitVec 32) (ink black : BitVec 8)
    (h : word.toNat = 257 * (255 - ink.toNat) * (255 - black.toNat) / 255) :
    ((word >>> 8).truncate 8).toNat =
      257 * (255 - ink.toNat) * (255 - black.toNat) / 65280 := by
  change (word >>> 8).toNat % 256 = _
  rw [BitVec.toNat_ushiftRight, h, Nat.shiftRight_eq_div_pow, Nat.div_div_eq_div_mul]
  change (257 * (255 - ink.toNat) * (255 - black.toNat) / 65280) % 256 = _
  apply Nat.mod_eq_of_lt
  apply (Nat.div_lt_iff_lt_mul (by decide)).2
  have product : 257 * (255 - ink.toNat) * (255 - black.toNat) ≤ 65535 * 255 :=
    Nat.mul_le_mul (by omega) (by omega)
  omega

/-- Each byte is the truncated result of the rational ink–black scale. -/
theorem cmykToRGB_spec (c m y k : BitVec 8) :
    (let (r, g, b) := cmykToRGB c m y k
     (r.toNat, g.toNat, b.toNat)) =
      (257 * (255 - c.toNat) * (255 - k.toNat) / 65280,
       257 * (255 - m.toNat) * (255 - k.toNat) / 65280,
       257 * (255 - y.toNat) * (255 - k.toNat) / 65280) := by
  have widened := cmykRGBA_spec c m y k
  generalize h : cmykRGBA c m y k = out at widened
  rcases out with ⟨r, g, b, a⟩
  simp only [Prod.mk.injEq] at widened
  rw [cmykToRGB, h]
  simp only [narrow_spec r c k widened.1, narrow_spec g m k widened.2.1,
    narrow_spec b y k widened.2.2.1]

end GoLean.Image.Color
