import GoLean.Slices.SwapRange

namespace GoLean.Slices

private def rotateSwap (length a b count : Nat) (cells : Option (List α)) :
    Except Unit (Option (List α)) :=
  let result := swapRangeOrdered cells length (Int.ofNat a) (Int.ofNat b) (Int.ofNat count)
  if result.boundsPanic then .error () else .ok result.cells

private def rotateBlocks (length middle : Nat) :
    Nat → Nat → Nat → Option (List α) → Except Unit (Option (List α))
  | 0, _, _, _ => .error ()
  | fuel + 1, i, j, cells =>
    if i = j then rotateSwap length (middle - i) middle i cells
    else if j < i then do
      let updated ← rotateSwap length (middle - i) middle j cells
      rotateBlocks length middle fuel (i - j) j updated
    else do
      let updated ← rotateSwap length (middle - i) (middle + j - i) i cells
      rotateBlocks length middle fuel i (j - i) updated

private def rotateRange (cells : Option (List α)) (length a middle b : Nat) :
    Except Unit (Option (List α)) :=
  if b ≤ length && (a = middle ∧ middle = b ∨ a < middle ∧ middle < b) then
    -- In the unequal-block loop one positive length strictly decreases.
    rotateBlocks length middle (b - a + 1) (middle - a) (b - middle) cells
  else .error ()

/-- Go 1.27.0 src/slices/zsortordered.go:466–481. Rotate two nonempty
    neighboring blocks (or two empty blocks) by the pinned swapping loop. -/
def rotateOrdered (cells : Option (List α)) (length a middle b : Nat) :
    Except Unit (Option (List α)) :=
  rotateRange cells length a middle b

/-- Go 1.27.0 src/slices/zsortanyfunc.go:464–479. Comparator is never used. -/
def rotateCmpFunc (cells : Option (List α)) (length a middle b : Nat)
    (_ : Option (α → α → Int)) : Except Unit (Option (List α)) :=
  rotateRange cells length a middle b

theorem rotateOrdered_empty (cells : Option (List α)) (length a : Nat) (bound : a ≤ length) :
    rotateOrdered cells length a a a = .ok cells := by
  simp [rotateOrdered, rotateRange, rotateBlocks, rotateSwap, bound,
    swapRangeOrdered_nonpositive]

theorem rotateOrdered_unequal :
    rotateOrdered (some [99, 1, 2, 3, 4, 5, 6, 7, 88, 77]) 9 1 4 8 =
      .ok (some [99, 4, 5, 6, 7, 1, 2, 3, 88, 77]) := rfl

theorem rotateCmpFunc_ignoresComparator (cells : Option (List α)) (length a middle b : Nat)
    (first second : Option (α → α → Int)) :
    rotateCmpFunc cells length a middle b first =
      rotateCmpFunc cells length a middle b second := rfl

theorem rotateCmpFunc_spec (cells : Option (List α)) (length a middle b : Nat)
    (compare : Option (α → α → Int)) :
    rotateCmpFunc cells length a middle b compare = rotateOrdered cells length a middle b := rfl

end GoLean.Slices
