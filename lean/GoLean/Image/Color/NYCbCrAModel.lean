import GoLean.Image.Color.YCbCrModel

namespace GoLean.Image.Color

abbrev NYCbCrAChannels := Int × Int × Int × Nat

inductive NYCbCrAInput where
  | nYCbCrA (channels : NYCbCrAChannels)
  | yCbCr (channels : YCbCrChannels)
  | other (callback : Option RGBAWords)
  | nilColor

private def unpremultiply (a x : BitVec 32) : BitVec 32 :=
  if a = 0 then x else (x * 65535) / a

private def highByte (x : BitVec 32) : BitVec 8 :=
  (x >>> 8).truncate 8

/-- A missing result denotes a panic. The count records RGBA calls.
    The multiplication in the nonzero-alpha path wraps before division.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/ycbcr.go#L294-L312 -/
def nYCbCrAModel : NYCbCrAInput → Option NYCbCrAChannels × Nat
  | .nYCbCrA channels => (some channels, 0)
  | .yCbCr (y, cb, cr) => (some (y, cb, cr, 255), 0)
  | .other callback =>
    (callback.map fun (r, g, b, a) =>
      let (y, cb, cr) := rgbToYCbCr (highByte (unpremultiply a r))
        (highByte (unpremultiply a g)) (highByte (unpremultiply a b))
      (y, cb, cr, (highByte a).toNat), 1)
  | .nilColor => (none, 0)

private theorem highByte_spec (w : BitVec 32) :
    highByte w = BitVec.ofNat 8 (w.toNat / 256) := by
  rw [← BitVec.toNat_inj]
  simp [highByte, BitVec.toNat_ushiftRight, Nat.shiftRight_eq_div_pow]

private theorem unpremultiply_spec (a x : BitVec 32) :
    (unpremultiply a x).toNat =
      if a.toNat = 0 then x.toNat else x.toNat * 65535 % 4294967296 / a.toNat := by
  have zero : a = 0 ↔ a.toNat = 0 := by rw [← BitVec.toNat_inj]; rfl
  unfold unpremultiply
  split
  · rename_i h
    rw [ite_eq_left (zero.mp h)]
  · rename_i h
    rw [ite_eq_right (fun n => h (zero.mpr n)), BitVec.toNat_udiv, BitVec.toNat_mul]
    rfl

theorem nYCbCrAModel_values (channels : NYCbCrAChannels) (y cb cr : Int) :
    nYCbCrAModel (.nYCbCrA channels) = (some channels, 0) ∧
    nYCbCrAModel (.yCbCr (y, cb, cr)) = (some (y, cb, cr, 255), 0) :=
  ⟨rfl, rfl⟩

theorem nYCbCrAModel_return (r g b a : BitVec 32) :
    nYCbCrAModel (.other (some (r, g, b, a))) =
      let byte := fun x : BitVec 32 => BitVec.ofNat 8
        ((if a.toNat = 0 then x.toNat else x.toNat * 65535 % 4294967296 / a.toNat) / 256)
      let (y, cb, cr) := rgbToYCbCr (byte r) (byte g) (byte b)
      (some (y, cb, cr, a.toNat / 256 % 256), 1) := by
  simp [nYCbCrAModel, highByte_spec, unpremultiply_spec]

theorem nYCbCrAModel_panic :
    nYCbCrAModel (.other none) = (none, 1) ∧ nYCbCrAModel .nilColor = (none, 0) :=
  ⟨rfl, rfl⟩

end GoLean.Image.Color
