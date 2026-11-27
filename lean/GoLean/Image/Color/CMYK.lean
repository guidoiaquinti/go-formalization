import Init.Data.BitVec.Lemmas

namespace GoLean.Image.Color

private def complementInk (ink : BitVec 8) : BitVec 32 :=
  0xffff - ink.zeroExtend 32 * 0x101

/-- The Go method forms the complement of widened K once, and scales each
    complemented ink channel by it before setting opaque alpha.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/ycbcr.go#L344-L361 -/
def cmykRGBA (c m y k : BitVec 8) : BitVec 32 × BitVec 32 × BitVec 32 × BitVec 32 :=
  let white := complementInk k
  (complementInk c * white / 0xffff, complementInk m * white / 0xffff,
   complementInk y * white / 0xffff, 0xffff)

private theorem complementInk_spec (ink : BitVec 8) :
    (complementInk ink).toNat = 257 * (255 - ink.toNat) := by
  have bound : ink.toNat < 256 := ink.isLt
  have scaled : (ink.zeroExtend 32 * (0x101 : BitVec 32)).toNat = ink.toNat * 257 := by
    rw [BitVec.toNat_mul, BitVec.toNat_setWidth_of_le (by decide : 8 ≤ 32)]
    change (ink.toNat * 257) % 2 ^ 32 = ink.toNat * 257
    exact Nat.mod_eq_of_lt (by omega)
  have fits : ink.zeroExtend 32 * (0x101 : BitVec 32) ≤ (0xffff : BitVec 32) := by
    change (ink.zeroExtend 32 * (0x101 : BitVec 32)).toNat ≤ 65535
    rw [scaled]
    omega
  rw [complementInk, BitVec.toNat_sub_of_le fits, scaled]
  change 65535 - ink.toNat * 257 = _
  omega

private theorem cmykChannel_spec (ink black : BitVec 8) :
    (complementInk ink * complementInk black / (0xffff : BitVec 32)).toNat =
      257 * (255 - ink.toNat) * (255 - black.toNat) / 255 := by
  have noWrap : (complementInk ink).toNat * (complementInk black).toNat < 2 ^ 32 := by
    rw [complementInk_spec, complementInk_spec]
    have bounded : 257 * (255 - ink.toNat) * (257 * (255 - black.toNat)) ≤ 65535 * 65535 :=
      Nat.mul_le_mul (by omega) (by omega)
    omega
  rw [BitVec.toNat_udiv, BitVec.toNat_mul, Nat.mod_eq_of_lt noWrap,
    complementInk_spec, complementInk_spec]
  change (257 * (255 - ink.toNat) * (257 * (255 - black.toNat))) / (257 * 255) = _
  rw [show 257 * (255 - ink.toNat) * (257 * (255 - black.toNat)) =
      257 * (257 * (255 - ink.toNat) * (255 - black.toNat)) by ac_rfl]
  exact Nat.mul_div_mul_left _ _ (by decide)

theorem cmykRGBA_spec (c m y k : BitVec 8) :
    (let (red, green, blue, alpha) := cmykRGBA c m y k
     (red.toNat, green.toNat, blue.toNat, alpha.toNat)) =
      (257 * (255 - c.toNat) * (255 - k.toNat) / 255,
       257 * (255 - m.toNat) * (255 - k.toNat) / 255,
       257 * (255 - y.toNat) * (255 - k.toNat) / 255,
       65535) := by
  simp only [cmykRGBA, cmykChannel_spec]
  rfl

end GoLean.Image.Color
