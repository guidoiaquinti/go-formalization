import GoLean.Image.Color.YCbCr

namespace GoLean.Image.Color

/-- The RGB computation repeats YCbCr.RGBA. Go widens alpha by replication
    and applies it with an unsigned multiplication and division.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/ycbcr.go#L247-L288 -/
def nYCbCrARGBA (y cb cr a : BitVec 8) : Int × Int × Int × Int :=
  let (r, g, b, _) := yCbCrRGBA y cb cr
  let alpha := (a.toNat : Int) * 257
  let applyAlpha := fun (c : Int) => (c * alpha % 4294967296) / 65535
  (applyAlpha r, applyAlpha g, applyAlpha b, alpha)

private theorem applyAlpha_spec (c : Int) (a : BitVec 8)
    (lo : 0 ≤ c) (hi : c ≤ 65535) :
    (c * ((a.toNat : Int) * 257) % 4294967296) / 65535 =
      c * (a.toNat : Int) / 255 := by
  have hA : (a.toNat : Int) * 257 ≤ 65535 := by have := a.isLt; omega
  have hZ : 0 ≤ (a.toNat : Int) * 257 := by omega
  have hMul := Int.mul_le_mul_of_nonneg_left hA lo
  have hBound := Int.mul_le_mul_of_nonneg_right hi (by decide : (0 : Int) ≤ 65535)
  rw [Int.emod_eq_of_lt (Int.mul_nonneg lo hZ) (by omega)]
  rw [← Int.mul_assoc]
  exact Int.mul_ediv_mul_of_pos_left (c * (a.toNat : Int)) 255 (by decide : (0 : Int) < 257)

theorem nYCbCrARGBA_spec (y cb cr a : BitVec 8) :
    let l := (y.toNat : Int)
    let u := (cb.toNat : Int) - 128
    let v := (cr.toNat : Int) - 128
    let scale := fun x : Int => min 65535 (max 0 x) * (a.toNat : Int) / 255
    nYCbCrARGBA y cb cr a =
      (scale (257 * l + (l + 91881 * v) / 256),
       scale (257 * l + (l - 22554 * u - 46802 * v) / 256),
       scale (257 * l + (l + 116130 * u) / 256), (a.toNat : Int) * 257) := by
  have scale (x : Int) := applyAlpha_spec (min 65535 (max 0 x)) a
    (Int.le_min.mpr ⟨by decide, Int.le_max_left _ _⟩) (Int.min_le_left _ _)
  simp only [nYCbCrARGBA, yCbCrRGBA_spec, scale]

end GoLean.Image.Color
