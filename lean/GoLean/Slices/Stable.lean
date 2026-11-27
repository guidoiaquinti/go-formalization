import GoLean.Slices.RotateRange
import GoLean.Slices.SortStateOps

namespace GoLean.Slices

private def search [Inhabited α]
    (goRight : LocalSortState α ε → Nat → Bool × LocalSortState α ε) :
    Nat → Nat → Nat → LocalSortState α ε → Nat × LocalSortState α ε
  | 0, lo, _, s => (lo, s)
  | fuel + 1, lo, hi, s =>
    if hi ≤ lo || s.failure.isSome then (lo, s) else
      let mid := (lo + hi) / 2
      let (right, s) := goRight s mid
      if right then search goRight fuel (mid + 1) hi s
      else search goRight fuel lo mid s

private def merge [Inhabited α] (compare : LocalCompare α ε) :
    Nat → Nat → Nat → Nat → LocalSortState α ε → LocalSortState α ε
  | 0, _, _, _, s => s
  | fuel + 1, a, m, b, s =>
    if s.failure.isSome then s
    else if m - a = 1 then
      let (i, s) := search (fun s h => lessAt compare s h a) (b - m) m b s
      (List.range (i - 1 - a)).foldl (fun s n => exchange s (a + n) (a + n + 1)) s
    else if b - m = 1 then
      let (i, s) := search (fun s h =>
        let (less, s) := lessAt compare s m h
        (!less, s)) (m - a) a m s
      (List.range (m - i)).reverse.foldl (fun s n => exchange s (i + n + 1) (i + n)) s
    else
      let mid := (a + b) / 2
      let n := mid + m
      let (lo, hi) := if mid < m then (n - b, mid) else (a, m)
      let (start, s) := search (fun s c =>
        let (less, s) := lessAt compare s (n - 1 - c) c
        (!less, s)) (hi - lo) lo hi s
      let stop := n - start
      if s.failure.isSome then s else
        let s := if start < m && m < stop then
            match rotateOrdered (some s.cells.toList) s.cells.size start m stop with
            | .ok (some cells) => { s with cells := cells.toArray }
            | _ => s
          else s
        let s := if a < start && start < mid then merge compare fuel a start mid s else s
        if mid < stop && stop < b then merge compare fuel mid stop b s else s

/-- Go 1.27.0 src/slices/zsortordered.go:380–460. The two runs are
    nonempty. The single-element paths use opposite equality rules. -/
def symMergeOrdered [Inhabited α] (s : LocalSortState α ε) (a m b : Nat)
    (less : α → α → Bool) : LocalSortState α ε :=
  merge (orderedCompare less) (b - a) a m b s

/-- Go 1.27.0 src/slices/zsortanyfunc.go:378–458. Comparator observations
    use the evolving backing and include the panic-producing call. -/
def symMergeCmpFunc [Inhabited α] (s : LocalSortState α ε) (a m b : Nat)
    (compare : Option (Nat → α → α → Except ε Int)) : LocalSortState α ε :=
  merge (scriptedCompare compare) (b - a) a m b s

private def mergePass [Inhabited α]
    (merge : LocalSortState α ε → Nat → Nat → Nat → LocalSortState α ε) (n width : Nat) :
    Nat → Nat → LocalSortState α ε → LocalSortState α ε
  | 0, _, s => s
  | fuel + 1, a, s =>
    let mid := a + width
    if n ≤ mid || s.failure.isSome then s
    else
      let b := a + 2 * width
      mergePass merge n width fuel b (merge s a mid (min b n))

private def mergePasses [Inhabited α]
    (merge : LocalSortState α ε → Nat → Nat → Nat → LocalSortState α ε) (n : Nat) :
    Nat → Nat → LocalSortState α ε → LocalSortState α ε
  | 0, _, s => s
  | fuel + 1, width, s =>
    if n ≤ width || s.failure.isSome then s
    else mergePasses merge n fuel (2 * width) (mergePass merge n width n 0 s)

private def stable [Inhabited α]
    (insert : LocalSortState α ε → Nat → Nat → LocalSortState α ε)
    (merge : LocalSortState α ε → Nat → Nat → Nat → LocalSortState α ε)
    (s : LocalSortState α ε) (n : Nat) : LocalSortState α ε :=
  let s := (List.range (n / 20)).foldl (fun s i => insert s (i * 20) ((i + 1) * 20)) s
  let s := insert s (n / 20 * 20) n
  mergePasses merge n n 20 s

/-- Go 1.27.0 src/slices/zsortordered.go:337–357. Sort blocks of twenty,
    then stably merge adjacent runs with doubling width. -/
def stableOrdered [Inhabited α] (s : LocalSortState α ε) (n : Nat)
    (less : α → α → Bool) : LocalSortState α ε :=
  stable (fun s a b => insertionSortOrdered s a b less)
    (fun s a m b => symMergeOrdered s a m b less) s n

/-- Go 1.27.0 src/slices/zsortanyfunc.go:335–355. Observes the same
    callback sequence across the initial blocks and every merge pass. -/
def stableCmpFunc [Inhabited α] (s : LocalSortState α ε) (n : Nat)
    (compare : Option (Nat → α → α → Except ε Int)) : LocalSortState α ε :=
  stable (fun s a b => insertionSortCmpFunc s a b compare)
    (fun s a m b => symMergeCmpFunc s a m b compare) s n

theorem symMergeOrdered_single :
    symMergeOrdered (⟨#[9, 4, 1, 3, 5, 8], none, []⟩ : LocalSortState Int Nat)
      1 2 5 (fun x y => x < y) =
      ⟨#[9, 1, 3, 4, 5, 8], none, []⟩ := rfl

theorem symMergeCmpFunc_rightLastTie :
    symMergeCmpFunc (⟨#[(2, 0), (2, 1), (2, 2)], none, []⟩ : LocalSortState (Int × Nat) Nat)
      0 2 3 (some (fun _ x y => .ok (x.1 - y.1))) =
      ⟨#[(2, 0), (2, 1), (2, 2)], none, [((2, 2), (2, 1))]⟩ := rfl

theorem stableOrdered_prefix :
    stableOrdered (⟨#[4, 1, 3, 2, 0], none, []⟩ : LocalSortState Int Nat)
      4 (fun x y => x < y) = ⟨#[1, 2, 3, 4, 0], none, []⟩ := rfl

theorem stableCmpFunc_partial :
    stableCmpFunc (⟨#[3, 2, 1, 0], none, []⟩ : LocalSortState Int Nat) 3
      (some (fun i _ _ => if i = 1 then .error 1 else .ok (-1))) =
      ⟨#[2, 3, 1, 0], some (.callback 1), [(2, 3), (1, 3)]⟩ := rfl

end GoLean.Slices
