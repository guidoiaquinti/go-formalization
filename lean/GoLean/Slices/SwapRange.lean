namespace GoLean.Slices

structure SwapRangeResult (α : Type) where
  cells : Option (List α)
  boundsPanic : Bool

private def swapRangeStep (length : Nat) (a b : Int) (i : Nat)
    (result : SwapRangeResult α) : SwapRangeResult α :=
  if result.boundsPanic then result
  else
    let x := a + Int.ofNat i
    let y := b + Int.ofNat i
    if x < 0 || y < 0 || length ≤ x.toNat || length ≤ y.toNat then
      {result with boundsPanic := true}
    else
      let cells := result.cells.getD []
      match cells[x.toNat]?, cells[y.toNat]? with
      | some left, some right =>
        ⟨some ((cells.set x.toNat right).set y.toNat left), false⟩
      | _, _ => {result with boundsPanic := true}

private def swapRange (cells : Option (List α)) (length : Nat) (a b count : Int) :
    SwapRangeResult α :=
  -- At most length swaps fit. The next access either panics or the loop
  -- has already finished, so finite storage bounds fuel even for MaxInt.
  (List.range (min count.toNat (length + 1))).foldl
    (fun result i => swapRangeStep length a b i result) ⟨cells, false⟩

/-- Go 1.27.0 src/slices/zsortordered.go:331–335. Each iteration reads its
    two current cells before writing them, so overlapping ranges compose
    the swaps in increasing index order. -/
def swapRangeOrdered (cells : Option (List α)) (length : Nat) (a b count : Int) :
    SwapRangeResult α :=
  swapRange cells length a b count

/-- Go 1.27.0 src/slices/zsortanyfunc.go:329–333. The comparator is unused. -/
def swapRangeCmpFunc (cells : Option (List α)) (length : Nat) (a b count : Int)
    (_ : Option (α → α → Int)) : SwapRangeResult α :=
  swapRange cells length a b count

theorem swapRangeOrdered_nonpositive (cells : Option (List α)) (length : Nat)
    (a b count : Int) (nonpositive : count ≤ 0) :
    swapRangeOrdered cells length a b count = ⟨cells, false⟩ := by
  simp [swapRangeOrdered, swapRange, show count.toNat = 0 by omega]

theorem swapRangeOrdered_overlap :
    swapRangeOrdered (some [10, 20, 30, 40, 99]) 4 0 1 3 =
      ⟨some [20, 30, 40, 10, 99], false⟩ := rfl

theorem swapRangeOrdered_partialPanic :
    swapRangeOrdered (some [10, 20, 30, 40, 99]) 4 0 2 3 =
      ⟨some [30, 40, 10, 20, 99], true⟩ := rfl

theorem swapRangeCmpFunc_ignoresComparator (cells : Option (List α)) (length : Nat)
    (a b count : Int) (first second : Option (α → α → Int)) :
    swapRangeCmpFunc cells length a b count first =
      swapRangeCmpFunc cells length a b count second := rfl

theorem swapRangeCmpFunc_spec (cells : Option (List α)) (length : Nat) (a b count : Int)
    (compare : Option (α → α → Int)) :
    swapRangeCmpFunc cells length a b count compare = swapRangeOrdered cells length a b count := rfl

end GoLean.Slices
