import GoLean.Image.Color.RGBToCMYK
import GoLean.Image.Color.YCbCrModel

namespace GoLean.Image.Color

abbrev CMYKChannels := BitVec 8 × BitVec 8 × BitVec 8 × BitVec 8

inductive CMYKInput where
  | cmyk (channels : CMYKChannels)
  | other (callback : Option RGBAWords)
  | nilColor

/-- A missing result denotes a panic. The count records RGBA calls.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/ycbcr.go#L366-L373 -/
def cmykModel : CMYKInput → Option CMYKChannels × Nat
  | .cmyk channels => (some channels, 0)
  | .other callback =>
    (callback.map fun (r, g, b, _) =>
      rgbToCMYK ((r >>> 8).truncate 8) ((g >>> 8).truncate 8) ((b >>> 8).truncate 8), 1)
  | .nilColor => (none, 0)

theorem cmykModel_value (channels : CMYKChannels) :
    cmykModel (.cmyk channels) = (some channels, 0) := rfl

private theorem cmykSecondByte (word : BitVec 32) :
    (word >>> 8).truncate 8 = BitVec.ofNat 8 (word.toNat / 256) := by
  rw [← BitVec.toNat_inj]
  simp [BitVec.toNat_ushiftRight, Nat.shiftRight_eq_div_pow]

theorem cmykModel_return (r g b a : BitVec 32) :
    cmykModel (.other (some (r, g, b, a))) =
      (some (rgbToCMYK (BitVec.ofNat 8 (r.toNat / 256))
        (BitVec.ofNat 8 (g.toNat / 256)) (BitVec.ofNat 8 (b.toNat / 256))), 1) := by
  simp [cmykModel, cmykSecondByte]

theorem cmykModel_panic : cmykModel (.other none) = (none, 1) := rfl

theorem cmykModel_nil : cmykModel .nilColor = (none, 0) := rfl

end GoLean.Image.Color
