import GoLean.Slices.SortPrimitives
import GoLean.Slices.SwapRange

namespace GoLean.Slices

private def breakPatterns (width : Nat) (cells : Option (List α)) (length : Nat)
    (a b : Int) : SwapRangeResult α :=
  let size := b - a
  let initial : SwapRangeResult α := ⟨cells, false⟩
  if size < 8 then initial
  else
    let mask := nextPowerOfTwo width size - 1
    let first := a + (size / 4) * 2 - 1
    let (_, result) := (List.range 3).foldl (fun (state, result) i =>
      if result.boundsPanic then (state, result)
      else
        match xorshiftNext state with
        | none => (state, result)
        | some (word, next) =>
          let candidate := Int.ofNat (BitVec.ofNat width word.toNat &&& mask).toNat
          let other := if candidate ≥ size then candidate - size else candidate
          (some next, swapRangeOrdered result.cells length (first + Int.ofNat i) (a + other) 1)
    ) (some (BitVec.ofInt 64 size), initial)
    result

/-- Go 1.27.0 src/slices/zsortordered.go:242–255. The three swaps are
    sequential. Native signed index arithmetic is assumed not to overflow. -/
def breakPatternsOrdered (width : Nat) (cells : Option (List α)) (length : Nat)
    (a b : Int) : SwapRangeResult α :=
  breakPatterns width cells length a b

/-- Go 1.27.0 src/slices/zsortanyfunc.go:240–253. The comparator is unused. -/
def breakPatternsCmpFunc (width : Nat) (cells : Option (List α)) (length : Nat)
    (a b : Int) (_ : Option (α → α → Int)) : SwapRangeResult α :=
  breakPatterns width cells length a b

theorem breakPatternsOrdered_short (width : Nat) (cells : Option (List α)) (length : Nat)
    (a b : Int) (short : b - a < 8) :
    breakPatternsOrdered width cells length a b = ⟨cells, false⟩ := by
  simp [breakPatternsOrdered, breakPatterns, short]

theorem breakPatternsOrdered_partialPanic :
    breakPatternsOrdered 64 (some [10, 20, 30, 40, 99]) 4 0 8 =
      ⟨some [40, 20, 30, 10, 99], true⟩ := rfl

theorem breakPatternsCmpFunc_ignoresComparator (width : Nat) (cells : Option (List α))
    (length : Nat) (a b : Int) (compare : Option (α → α → Int)) :
    breakPatternsCmpFunc width cells length a b compare =
      breakPatternsOrdered width cells length a b := rfl

end GoLean.Slices
