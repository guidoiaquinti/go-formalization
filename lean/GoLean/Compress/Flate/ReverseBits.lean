import GoLean.Math.Bits.Reversal

namespace GoLean.Compress.Flate

/-- Go wraps the byte subtraction and masks its low nibble before shifting. -/
def reverseBits (x : BitVec 16) (count : BitVec 8) : BitVec 16 :=
  GoLean.Math.Bits.reverse16 (x <<< (((16 : BitVec 8) - count) &&& 15).toNat)

theorem reverseBits_bit (x : BitVec 16) (count : BitVec 8) (i : Nat) (hi : i < 16) :
    (reverseBits x count).getLsbD i =
      let shift := (((16 : BitVec 8) - count) &&& 15).toNat
      !decide (15 - i < shift) && x.getLsbD (15 - i - shift) := by
  unfold reverseBits
  rw [GoLean.Math.Bits.reverse16_spec _ i hi, BitVec.getLsbD_shiftLeft]
  have h : 15 - i < 16 := by omega
  simp [h]

end GoLean.Compress.Flate
