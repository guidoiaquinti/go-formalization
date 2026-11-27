import GoLean.Image.Color.RGBToYCbCr

namespace GoLean.Image.Color

abbrev YCbCrChannels := Int × Int × Int
abbrev RGBAWords := BitVec 32 × BitVec 32 × BitVec 32 × BitVec 32

inductive YCbCrInput where
  | yCbCr (channels : YCbCrChannels)
  | other (callback : Option RGBAWords)
  | nilColor

/-- A missing result denotes a panic. The count records RGBA calls.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/ycbcr.go#L231-L238 -/
def yCbCrModel : YCbCrInput → Option YCbCrChannels × Nat
  | .yCbCr channels => (some channels, 0)
  | .other callback =>
    (callback.map fun (r, g, b, _) =>
      rgbToYCbCr ((r >>> 8).truncate 8) ((g >>> 8).truncate 8) ((b >>> 8).truncate 8), 1)
  | .nilColor => (none, 0)

theorem yCbCrModel_value (channels : YCbCrChannels) :
    yCbCrModel (.yCbCr channels) = (some channels, 0) := rfl

private theorem secondByte (word : BitVec 32) :
    (word >>> 8).truncate 8 = BitVec.ofNat 8 (word.toNat / 256) := by
  rw [← BitVec.toNat_inj]
  simp [BitVec.toNat_ushiftRight, Nat.shiftRight_eq_div_pow]

theorem yCbCrModel_return (r g b a : BitVec 32) :
    yCbCrModel (.other (some (r, g, b, a))) =
      (some (rgbToYCbCr (BitVec.ofNat 8 (r.toNat / 256))
        (BitVec.ofNat 8 (g.toNat / 256)) (BitVec.ofNat 8 (b.toNat / 256))), 1) := by
  simp [yCbCrModel, secondByte]

theorem yCbCrModel_panic : yCbCrModel (.other none) = (none, 1) := rfl

theorem yCbCrModel_nil : yCbCrModel .nilColor = (none, 0) := rfl

end GoLean.Image.Color
