import Init.Data.BitVec.Lemmas

namespace GoLean.Image.Color

/-- The source tests the high byte of the unsigned reinterpretation, then
    uses a signed right shift, complement, and conversion to uint8.
    Here division by a power of two represents arithmetic right shift,
    remainder represents unsigned narrowing, and complement is -x-1. -/
private def channel (x : Int) : Int :=
  (if (x % 4294967296) / 16777216 = 0 then x / 65536 else -(x / 2147483648) - 1) % 256

/-- Go 1.27.0, src/image/color/ycbcr.go, lines 115–153.
    The products and intermediate sums are within signed 32-bit range.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/ycbcr.go#L115-L153 -/
def yCbCrToRGB (y cb cr : BitVec 8) : Int × Int × Int :=
  let yy := (y.toNat : Int) * 65793
  let u := (cb.toNat : Int) - 128
  let v := (cr.toNat : Int) - 128
  (channel (yy + 91881 * v),
    channel (yy - 22554 * u - 46802 * v),
    channel (yy + 116130 * u))

private theorem channel_eq_clamp (x : Int) (lo : -2147483648 ≤ x) (hi : x < 2147483648) :
    channel x = min 255 (max 0 (x / 65536)) := by
  unfold channel
  split <;> simp only [Int.min_def, Int.max_def] <;> split <;> split <;> omega

/-- The result adds luminance to the rounded chroma adjustment and clamps
    each channel to the byte interval. -/
theorem yCbCrToRGB_spec (y cb cr : BitVec 8) :
    let l := (y.toNat : Int)
    let u := (cb.toNat : Int) - 128
    let v := (cr.toNat : Int) - 128
    yCbCrToRGB y cb cr =
      (min 255 (max 0 (l + (257 * l + 91881 * v) / 65536)),
        min 255 (max 0 (l + (257 * l - 22554 * u - 46802 * v) / 65536)),
        min 255 (max 0 (l + (257 * l + 116130 * u) / 65536))) := by
  have hy : y.toNat < 256 := y.isLt
  have hu : cb.toNat < 256 := cb.isLt
  have hv : cr.toNat < 256 := cr.isLt
  simp only [yCbCrToRGB]
  rw [channel_eq_clamp _ (by omega) (by omega),
    channel_eq_clamp _ (by omega) (by omega),
    channel_eq_clamp _ (by omega) (by omega)]
  congr 3 <;> omega

end GoLean.Image.Color
