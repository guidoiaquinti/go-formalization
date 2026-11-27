import GoLean.Slices.Median

namespace GoLean.Slices

structure ChoosePivotObservation (α ε : Type) where
  result : Except (Order2Failure ε) (Nat × Nat)
  calls : List (α × α)

private structure PivotTrace (α : Type) where
  count : Option Int
  calls : List (α × α)

private def runMedian (step : Nat → Option Int → MedianObservation α ε) :
    ExceptT (Order2Failure ε) (StateM (PivotTrace α)) Nat :=
  ExceptT.mk fun state =>
    let m := step state.calls.length state.count
    (m.result, ⟨m.counter, state.calls ++ m.calls⟩)

private def choosePivot (a b : Nat)
    (adjacent : Nat → Nat → Option Int → MedianObservation α ε)
    (median : Nat → Nat → Nat → Nat → Option Int → MedianObservation α ε) :
    ChoosePivotObservation α ε :=
  let size := b - a
  let quarter := size / 4
  let i := a + quarter
  let j := a + quarter * 2
  let k := a + quarter * 3
  let work : ExceptT (Order2Failure ε) (StateM (PivotTrace α)) (Nat × Nat) := do
    if size < 8 then
      return (j, 1)
    let (i, j, k) ← if size ≥ 50 then do
      let i ← runMedian (adjacent i)
      let j ← runMedian (adjacent j)
      let k ← runMedian (adjacent k)
      pure (i, j, k)
    else pure (i, j, k)
    let pivot ← runMedian (median i j k)
    let trace ← get
    let hint := if trace.count = some 0 then 1 else if trace.count = some 12 then 2 else 0
    return (pivot, hint)
  let (result, trace) := work.run ⟨some 0, []⟩
  ⟨result, trace.calls⟩

/-- Go 1.27.0 src/slices/zsortordered.go:263–295, with a valid range. -/
def choosePivotOrdered (width : Nat) (get : Nat → α) (a b : Nat)
    (less : α → α → Bool) : ChoosePivotObservation α ε :=
  choosePivot a b
    (fun i _ count => medianAdjacentOrdered width get i count less)
    (fun i j k _ count => medianOrdered width get i j k count less)

/-- Go 1.27.0 src/slices/zsortanyfunc.go:261–293, with a valid range and
    comparator script indexed by the call number across all medians. -/
def choosePivotCmpFunc (width : Nat) (get : Nat → α) (a b : Nat)
    (compare : Option (Nat → α → α → Except ε Int)) : ChoosePivotObservation α ε :=
  let atCall (offset : Nat) := compare.map fun f i => f (offset + i)
  choosePivot a b
    (fun i offset count => medianAdjacentCmpFunc width get i count (atCall offset))
    (fun i j k offset count => medianCmpFunc width get i j k count (atCall offset))

theorem choosePivotOrdered_short (width a b : Nat) (values : Nat → α) (less : α → α → Bool)
    (short : b - a < 8) :
    choosePivotOrdered (ε := ε) width values a b less =
      ⟨.ok (a + ((b - a) / 4) * 2, 1), []⟩ := by
  simp [choosePivotOrdered, choosePivot, short]
  rfl

theorem choosePivotCmpFunc_short (width a b : Nat) (values : Nat → α)
    (compare : Option (Nat → α → α → Except ε Int)) (short : b - a < 8) :
    choosePivotCmpFunc width values a b compare =
      ⟨.ok (a + ((b - a) / 4) * 2, 1), []⟩ := by
  simp [choosePivotCmpFunc, choosePivot, short]
  rfl

theorem choosePivotOrdered_decreasing :
    choosePivotOrdered (ε := Unit) 64 (fun i => 50 - i) 0 50 (fun x y => x < y) =
      ⟨.ok (24, 2), []⟩ := rfl

theorem choosePivotCmpFunc_nintherTrace :
    (choosePivotCmpFunc 64 id 0 50 (some fun i _ _ =>
        if i = 3 then .error 3 else .ok 0)).calls =
      [(12, 11), (13, 12), (12, 11), (24, 23)] := rfl

end GoLean.Slices
