import Init.Data.BitVec.Lemmas

namespace GoLean.Image.Color

private def rgbInk (white channel : BitVec 32) : BitVec 8 :=
  ((white - channel) * 255 / white).truncate 8

private def selectWhite (a b : BitVec 32) : BitVec 32 :=
  if a < b then b else a

/-- Go widens the inputs, selects the brightest channel, and handles black before dividing.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/ycbcr.go#L314-L333 -/
def rgbToCMYK (r g b : BitVec 8) : BitVec 8 × BitVec 8 × BitVec 8 × BitVec 8 :=
  let rr := r.zeroExtend 32
  let gg := g.zeroExtend 32
  let bb := b.zeroExtend 32
  let white := selectWhite (selectWhite rr gg) bb
  if white == 0 then (0, 0, 0, 255)
  else (rgbInk white rr, rgbInk white gg, rgbInk white bb, (255 - white).truncate 8)

private theorem selectWhite_spec (a b : BitVec 32) :
    (selectWhite a b).toNat = max a.toNat b.toNat := by
  unfold selectWhite
  split
  next h => exact (Nat.max_eq_right (Nat.le_of_lt h)).symm
  next h => exact (Nat.max_eq_left (Nat.le_of_not_gt h)).symm

private theorem rgbInk_spec (white channel : BitVec 32)
    (positive : 0 < white.toNat) (le : channel.toNat ≤ white.toNat)
    (bound : white.toNat ≤ 255) :
    (rgbInk white channel).toNat = (white.toNat - channel.toNat) * 255 / white.toNat := by
  have h : channel ≤ white := le
  have product : (white.toNat - channel.toNat) * 255 ≤ white.toNat * 255 :=
    Nat.mul_le_mul_right 255 (Nat.sub_le _ _)
  have wraps : (white.toNat - channel.toNat) * 255 < 2 ^ 32 := by omega
  have result : (white.toNat - channel.toNat) * 255 / white.toNat < 256 := by
    apply (Nat.div_lt_iff_lt_mul positive).2
    omega
  change (((white - channel) * 255 / white).toNat % 256) = _
  rw [BitVec.toNat_udiv, BitVec.toNat_mul, BitVec.toNat_sub_of_le h]
  change (((white.toNat - channel.toNat) * 255 % 2 ^ 32) / white.toNat) % 256 = _
  rw [Nat.mod_eq_of_lt wraps, Nat.mod_eq_of_lt result]

theorem rgbToCMYK_spec (r g b : BitVec 8) :
    (let (c, m, y, k) := rgbToCMYK r g b
     (c.toNat, m.toNat, y.toNat, k.toNat)) =
      let white := max (max r.toNat g.toNat) b.toNat
      if white = 0 then (0, 0, 0, 255)
      else ((white - r.toNat) * 255 / white,
            (white - g.toNat) * 255 / white,
            (white - b.toNat) * 255 / white,
            255 - white) := by
  let w := selectWhite (selectWhite (r.zeroExtend 32) (g.zeroExtend 32)) (b.zeroExtend 32)
  have widen (x : BitVec 8) : (x.zeroExtend 32).toNat = x.toNat :=
    BitVec.toNat_setWidth_of_le (by decide)
  have hw : w.toNat = max (max r.toNat g.toNat) b.toNat := by
    simp only [w, selectWhite_spec, widen]
  have bound : w.toNat ≤ 255 := by
    have hr : r.toNat < 256 := r.isLt
    have hg : g.toNat < 256 := g.isLt
    have hb : b.toNat < 256 := b.isLt
    omega
  have leR : (r.zeroExtend 32).toNat ≤ w.toNat := by rw [widen, hw]; omega
  have leG : (g.zeroExtend 32).toNat ≤ w.toNat := by rw [widen, hw]; omega
  have leB : (b.zeroExtend 32).toNat ≤ w.toNat := by rw [widen, hw]; omega
  unfold rgbToCMYK
  dsimp only
  change (let (c, m, y, k) :=
            if w == 0 then (0, 0, 0, 255)
            else (rgbInk w (r.zeroExtend 32), rgbInk w (g.zeroExtend 32),
                  rgbInk w (b.zeroExtend 32), (255 - w).truncate 8)
          (c.toNat, m.toNat, y.toNat, k.toNat)) = _
  rw [← hw]
  by_cases h : w = 0
  · simp [h]
  · have positive : 0 < w.toNat := by
      have nz : w.toNat ≠ 0 := fun hz => h (BitVec.eq_of_toNat_eq hz)
      omega
    have black : (((255 : BitVec 32) - w).truncate 8).toNat = 255 - w.toNat := by
      have fits : w ≤ (255 : BitVec 32) := bound
      change ((255 - w).toNat % 256) = _
      rw [BitVec.toNat_sub_of_le fits]
      change (255 - w.toNat) % 256 = _
      exact Nat.mod_eq_of_lt (by omega)
    simp only [beq_iff_eq, h, ite_false, Nat.ne_of_gt positive,
      rgbInk_spec w _ positive leR bound, rgbInk_spec w _ positive leG bound,
      rgbInk_spec w _ positive leB bound, black, widen]

end GoLean.Image.Color
