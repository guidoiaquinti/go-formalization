namespace GoLean.Sort

/-- The callback arguments for reversing a valid half-open range. -/
def reverseRange (a b : Nat) : List (Nat × Nat) :=
  (List.range ((b - a) / 2)).map fun offset => (a + offset, b - 1 - offset)

/-- The callbacks run from low to high offsets, including for overlapping ranges. -/
def swapRange (a b n : Nat) : List (Nat × Nat) :=
  (List.range n).map fun offset => (a + offset, b + offset)

private def rotateSteps (m : Nat) : Nat → Nat → Nat → List (Nat × Nat)
  | 0, _, _ => []
  | fuel + 1, i, j =>
    if i = j then swapRange (m - i) m i
    else if j < i then
      swapRange (m - i) m j ++ rotateSteps m fuel (i - j) j
    else
      swapRange (m - i) (m + j - i) i ++ rotateSteps m fuel i (j - i)

/-- Each unequal-block iteration reduces the remaining positive total, so the original
range width gives enough steps. A range with two empty blocks emits no swaps. -/
def rotate (a m b : Nat) : List (Nat × Nat) :=
  rotateSteps m (b - a) (m - a) (b - m)

def reverseRangeFunc := reverseRange
def swapRangeFunc := swapRange
def rotateFunc := rotate

theorem reverseRange_length (a b : Nat) : (reverseRange a b).length = (b - a) / 2 := by
  simp [reverseRange]

theorem swapRange_length (a b n : Nat) : (swapRange a b n).length = n := by
  simp [swapRange]

theorem rotate_width_bound (a m b : Nat) (ha : a ≤ m) (hb : m ≤ b) :
    (rotate a m b).length ≤ b - a := by
  have steps (fuel i j : Nat) : (rotateSteps m fuel i j).length ≤ i + j := by
    induction fuel generalizing i j with
    | zero => simp [rotateSteps]
    | succ fuel ih =>
      unfold rotateSteps
      split
      · rw [swapRange_length]
        omega
      · split
        · rw [List.length_append, swapRange_length]
          have h := ih (i - j) j
          omega
        · rw [List.length_append, swapRange_length]
          have h := ih i (j - i)
          omega
  have h := steps (b - a) (m - a) (b - m)
  unfold rotate
  omega

theorem reverseRangeFunc_same (a b : Nat) : reverseRangeFunc a b = reverseRange a b := rfl
theorem swapRangeFunc_same (a b n : Nat) : swapRangeFunc a b n = swapRange a b n := rfl
theorem rotateFunc_same (a m b : Nat) : rotateFunc a m b = rotate a m b := rfl

end GoLean.Sort
