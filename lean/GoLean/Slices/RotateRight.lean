import GoLean.Slices.RotateLeft
import Init.Data.BitVec.Lemmas

namespace GoLean.Slices

/-- Go 1.27.0 src/slices/slices.go:450–452. The subtraction has Go's
    native signed int width, even for counts outside the visible length. -/
def sliceRotateRight (cells : List α) (length : Nat) (width : Nat) (r : Int) :
    Bool × List α :=
  sliceRotateLeft cells length (BitVec.ofInt width ((length : Int) - r)).toInt

theorem sliceRotateRight_spec (cells : List α) (length width r : Nat)
    (validLength : length ≤ cells.length) (widthPositive : 0 < width)
    (lengthFits : (length : Int) < 2 ^ (width - 1)) (validRotation : r ≤ length) :
    sliceRotateRight cells length width r =
      (true, (cells.take length).drop (length - r) ++ cells.take (length - r) ++ cells.drop length) := by
  have lower : -2 ^ (width - 1) ≤ (length : Int) - r := by
    have positive : 0 < (2 : Int) ^ (width - 1) := Int.pow_pos (by decide)
    omega
  have upper : (length : Int) - r < 2 ^ (width - 1) := by omega
  have difference : ((length : Int) - r) = ((length - r : Nat) : Int) := by omega
  unfold sliceRotateRight
  rw [BitVec.toInt_ofInt_eq_self widthPositive lower upper, difference]
  exact sliceRotateLeft_spec cells length (length - r) validLength (Nat.sub_le ..)

theorem sliceRotateRight_firstPanic (cells : List α) (length width : Nat) (r : Int)
    (outOfCapacity :
      ((length : Int) - r).bmod (2 ^ width) < 0 ∨
      cells.length < (((length : Int) - r).bmod (2 ^ width)).toNat) :
    sliceRotateRight cells length width r = (false, cells) := by
  rw [sliceRotateRight, BitVec.toInt_ofInt]
  exact sliceRotateLeft_firstPanic _ _ _ outOfCapacity

theorem sliceRotateRight_secondPanic (cells : List α) (length width leftCount : Nat) (r : Int)
    (adjusted : ((length : Int) - r).bmod (2 ^ width) = leftCount)
    (fitsCapacity : leftCount ≤ cells.length) (pastLength : length < leftCount) :
    sliceRotateRight cells length width r =
      (false, (cells.take leftCount).reverse ++ cells.drop leftCount) := by
  rw [sliceRotateRight, BitVec.toInt_ofInt, adjusted]
  exact sliceRotateLeft_secondPanic cells length leftCount fitsCapacity pastLength

end GoLean.Slices
