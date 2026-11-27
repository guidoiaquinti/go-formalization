import Init.Data.BitVec.Lemmas
import Init.Omega

namespace GoLean.Compress.Flate

/-- Starts of the DEFLATE length intervals, relative to match length three. -/
private def lengthStarts : List Nat :=
  [0, 1, 2, 3, 4, 5, 6, 7, 8, 10, 12, 14, 16, 20, 24, 28, 32, 40, 48,
   56, 64, 80, 96, 112, 128, 160, 192, 224, 255]

def lengthCode (length : BitVec 8) : Nat :=
  (lengthStarts.takeWhile (· ≤ length.toNat)).length - 1

private def shortCode (offset : Nat) : Nat :=
  if offset < 4 then offset
  else
    let high := offset.log2
    2 * high + (offset / 2 ^ (high - 1)) % 2

/-- The long branch drops seven low bits and narrows the result to a byte. -/
def offsetCode (offset : Nat) : Nat :=
  if offset < 256 then shortCode offset
  else 14 + shortCode (offset / 128 % 256)

theorem lengthCode_max : lengthCode (255 : BitVec 8) = 28 := by decide

theorem lengthCode_firstPaired :
    lengthCode (8 : BitVec 8) = 8 ∧ lengthCode (9 : BitVec 8) = 8 := by decide

theorem offsetCode_switch : offsetCode 255 = 15 ∧ offsetCode 256 = 16 := by decide

theorem offsetCode_periodic (offset : Nat) (long : 256 ≤ offset) :
    offsetCode (offset + 32768) = offsetCode offset := by
  have big : ¬ offset + 32768 < 256 := by omega
  have notShort : ¬ offset < 256 := by omega
  have same : (offset + 32768) / 128 % 256 = offset / 128 % 256 := by omega
  simp [offsetCode, big, notShort, same]

end GoLean.Compress.Flate
