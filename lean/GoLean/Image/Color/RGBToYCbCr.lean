import Init.Data.BitVec.Lemmas

namespace GoLean.Image.Color

/-- Go tests the high byte of the unsigned reinterpretation. Its arithmetic
    shifts are division by a power of two, the signed complement is -x-1,
    and converting to uint8 keeps the residue below 256. -/
private def chroma (x : Int) : Int :=
  (if (x % 4294967296) / 16777216 = 0 then x / 65536 else -(x / 2147483648) - 1) % 256

/-- The source widens the bytes to int32 before taking the weighted sums.
    All intermediate products and sums fit in int32.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/ycbcr.go#L7-L55 -/
def rgbToYCbCr (r g b : BitVec 8) : Int × Int × Int :=
  let r1 : Int := r.toNat
  let g1 : Int := g.toNat
  let b1 : Int := b.toNat
  (((19595 * r1 + 38470 * g1 + 7471 * b1 + 32768) / 65536) % 256,
    chroma (-11056 * r1 - 21712 * g1 + 32768 * b1 + 8421376),
    chroma (32768 * r1 - 27440 * g1 - 5328 * b1 + 8421376))

private theorem chroma_eq_clamp (x : Int) (lo : -2147483648 ≤ x) (hi : x < 2147483648) :
    chroma x = min 255 (max 0 (x / 65536)) := by
  unfold chroma
  split <;> simp only [Int.min_def, Int.max_def] <;> split <;> split <;> omega

/-- Express luminance and both chroma adjustments relative to green. -/
theorem rgbToYCbCr_spec (r g b : BitVec 8) :
    let red : Int := r.toNat
    let green : Int := g.toNat
    let blue : Int := b.toNat
    let dr := red - green
    let db := blue - green
    rgbToYCbCr r g b =
      (green + (19595 * dr + 7471 * db + 32768) / 65536,
        min 255 (max 0 (128 + (-11056 * dr + 32768 * db + 32768) / 65536)),
        min 255 (max 0 (128 + (32768 * dr - 5328 * db + 32768) / 65536))) := by
  have hr : r.toNat < 256 := r.isLt
  have hg : g.toNat < 256 := g.isLt
  have hb : b.toNat < 256 := b.isLt
  simp only [rgbToYCbCr]
  rw [chroma_eq_clamp _ (by omega) (by omega), chroma_eq_clamp _ (by omega) (by omega)]
  let yy := (19595 * (r.toNat : Int) + 38470 * (g.toNat : Int) +
    7471 * (b.toNat : Int) + 32768) / 65536
  have lo : 0 ≤ yy := by dsimp [yy]; omega
  have hi : yy < 256 := by dsimp [yy]; omega
  rw [Int.emod_eq_of_lt lo hi]
  congr 3 <;> omega

end GoLean.Image.Color
