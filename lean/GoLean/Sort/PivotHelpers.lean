import Init.Data.BitVec.Basic

namespace GoLean.Sort

structure PivotState where
  replies : List Bool
  calls : List (Nat × Nat)
  swaps : Int

private def stepLess (state : PivotState) (a b : Nat) : Bool × PivotState :=
  (state.replies.headD false, {state with replies := state.replies.tail, calls := (a, b) :: state.calls})

/-- Go 1.27.0 sort/zsortinterface.go:298–305. Reorders indices without
    writing the collection, while counting an inverted comparison. -/
def order2 (state : PivotState) (a b : Nat) : (Nat × Nat) × PivotState :=
  let (less, state) := stepLess state b a
  if less then ((b, a), {state with swaps := (BitVec.ofInt 64 (state.swaps + 1)).toInt})
  else ((a, b), state)

def order2Func := order2

/-- Three conditional index exchanges over a, b, c. -/
def median (state : PivotState) (a b c : Nat) : Nat × PivotState :=
  let ((a, b), state) := order2 state a b
  let ((b, _), state) := order2 state b c
  let ((_, b), state) := order2 state a b
  (b, state)

def medianFunc := median

def medianAdjacent (state : PivotState) (a : Nat) : Nat × PivotState :=
  median state (a - 1) a (a + 1)

def medianAdjacentFunc := medianAdjacent

/-- Go 1.27.0 sort/zsortinterface.go:261–295. Returns the pivot index
    and the numeric sorted hint, with exact comparator order. -/
def choosePivot (replies : List Bool) (a b : Nat) : (Nat × Nat) × PivotState :=
  let width := b - a
  let quarter := width / 4
  let i := a + quarter
  let j := a + 2 * quarter
  let k := a + 3 * quarter
  let state : PivotState := ⟨replies, [], 0⟩
  let (j, state) :=
    if width < 8 then (j, state)
    else if width < 50 then median state i j k
    else
      let (i, state) := medianAdjacent state i
      let (j, state) := medianAdjacent state j
      let (k, state) := medianAdjacent state k
      median state i j k
  let hint := if state.swaps = 0 then 1 else if state.swaps = 12 then 2 else 0
  ((j, hint), state)

def choosePivotFunc := choosePivot

theorem order2_tie (a b : Nat) (rest : List Bool) (swaps : Int) :
    order2 ⟨false :: rest, [], swaps⟩ a b = ((a, b), ⟨rest, [(b, a)], swaps⟩) := rfl

theorem order2_wraps :
    (order2 ⟨[true], [], 9223372036854775807⟩ 0 1).2.swaps = -9223372036854775808 := rfl

theorem order2Func_same (state : PivotState) (a b : Nat) :
    order2Func state a b = order2 state a b := rfl

theorem median_ties (a b c : Nat) (swaps : Int) :
    median ⟨[false, false, false], [], swaps⟩ a b c =
      (b, ⟨[], [(b, a), (c, b), (b, a)], swaps⟩) := rfl

theorem medianFunc_same (state : PivotState) (a b c : Nat) :
    medianFunc state a b c = median state a b c := rfl

theorem medianAdjacent_unfold (state : PivotState) (a : Nat) :
    medianAdjacent state a = median state (a - 1) a (a + 1) := rfl

theorem medianAdjacentFunc_same (state : PivotState) (a : Nat) :
    medianAdjacentFunc state a = medianAdjacent state a := rfl

theorem choosePivot_short (replies : List Bool) (a b : Nat) (short : b - a < 8) :
    choosePivot replies a b = ((a + 2 * ((b - a) / 4), 1), ⟨replies, [], 0⟩) := by
  simp [choosePivot, short]

theorem choosePivot_nintherDecreasing :
    (choosePivot (List.replicate 12 true) 0 50).1.2 = 2 := rfl

theorem choosePivotFunc_same (replies : List Bool) (a b : Nat) :
    choosePivotFunc replies a b = choosePivot replies a b := rfl

end GoLean.Sort
