import Init.Data.BitVec.Lemmas

namespace GoLean.Image.Color

/-- The high byte selects the shift or signed saturation path. Euclidean
    division implements arithmetic shifts and remainder implements the masks. -/
private def rgbaChannel (x : Int) : Int :=
  if (x % 4294967296) / 16777216 = 0 then x / 256
  else (-(x / 2147483648) - 1) % 65536

/-- Go keeps sixteen channel bits and sets alpha to full opacity.
    The signed products and sums fit in int32.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/ycbcr.go#L189-L225 -/
def yCbCrRGBA (y cb cr : BitVec 8) : Int × Int × Int × Int :=
  let yy := (y.toNat : Int) * 65793
  let u := (cb.toNat : Int) - 128
  let v := (cr.toNat : Int) - 128
  (rgbaChannel (yy + 91881 * v),
    rgbaChannel (yy - 22554 * u - 46802 * v),
    rgbaChannel (yy + 116130 * u), 65535)

private theorem rgbaChannel_eq_clamp (x : Int) (lo : -2147483648 ≤ x)
    (hi : x < 2147483648) :
    rgbaChannel x = min 65535 (max 0 (x / 256)) := by
  unfold rgbaChannel
  split <;> simp only [Int.min_def, Int.max_def] <;> split <;> split <;> omega

/-- Each channel adds a floor-rounded chroma adjustment to replicated
    luminance, then clamps to the inclusive sixteen-bit interval. -/
theorem yCbCrRGBA_spec (y cb cr : BitVec 8) :
    let l := (y.toNat : Int)
    let u := (cb.toNat : Int) - 128
    let v := (cr.toNat : Int) - 128
    yCbCrRGBA y cb cr =
      (min 65535 (max 0 (257 * l + (l + 91881 * v) / 256)),
        min 65535 (max 0 (257 * l + (l - 22554 * u - 46802 * v) / 256)),
        min 65535 (max 0 (257 * l + (l + 116130 * u) / 256)), 65535) := by
  have hy : y.toNat < 256 := y.isLt
  have hu : cb.toNat < 256 := cb.isLt
  have hv : cr.toNat < 256 := cr.isLt
  simp only [yCbCrRGBA]
  rw [rgbaChannel_eq_clamp _ (by omega) (by omega),
    rgbaChannel_eq_clamp _ (by omega) (by omega),
    rgbaChannel_eq_clamp _ (by omega) (by omega)]
  congr 3 <;> omega

end GoLean.Image.Color
