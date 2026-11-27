import GoLean.Slices.Reverse

namespace GoLean.Slices

/-- Go 1.27.0 src/slices/slices.go:443–449. The list includes all capacity
    cells. The first reverse can reach hidden cells before the second slice
    expression panics if the rotation exceeds the visible length. -/
def sliceRotateLeft (cells : List α) (length : Nat) (r : Int) : Bool × List α :=
  if r < 0 ∨ cells.length < r.toNat then (false, cells)
  else
    let first := sliceReverse (cells.take r.toNat)
    if length < r.toNat then (false, first ++ cells.drop r.toNat)
    else
      let second := sliceReverse ((cells.drop r.toNat).take (length - r.toNat))
      (true, sliceReverse (first ++ second) ++ cells.drop length)

theorem sliceRotateLeft_spec (cells : List α) (length r : Nat)
    (validLength : length ≤ cells.length) (validRotation : r ≤ length) :
    sliceRotateLeft cells length r =
      (true, (cells.take length).drop r ++ cells.take r ++ cells.drop length) := by
  have capacity : ¬ cells.length < r := by omega
  have inBounds : ¬ length < r := by omega
  simp [sliceRotateLeft, capacity, inBounds, sliceReverse_spec, List.take_drop,
    Nat.add_sub_cancel' validRotation, List.append_assoc]

theorem sliceRotateLeft_firstPanic (cells : List α) (length : Nat) (r : Int)
    (outOfCapacity : r < 0 ∨ cells.length < r.toNat) :
    sliceRotateLeft cells length r = (false, cells) := by
  rw [sliceRotateLeft, ite_eq_left outOfCapacity]

theorem sliceRotateLeft_secondPanic (cells : List α) (length r : Nat)
    (fitsCapacity : r ≤ cells.length) (pastLength : length < r) :
    sliceRotateLeft cells length r = (false, (cells.take r).reverse ++ cells.drop r) := by
  have nonnegative : ¬ (r : Int) < 0 := by omega
  simp [sliceRotateLeft, nonnegative, Nat.not_lt_of_ge fitsCapacity, pastLength, sliceReverse_spec]

end GoLean.Slices
