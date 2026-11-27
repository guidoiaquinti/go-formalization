namespace GoLean.Sort

structure LocalSortState where
  values : Array Int
  reversedCalls : List (Bool × Nat × Nat)

private def lessAt (state : LocalSortState) (i j : Nat) : Bool × LocalSortState :=
  (state.values[i]! < state.values[j]!,
    { state with reversedCalls := (false, i, j) :: state.reversedCalls })

private def exchange (state : LocalSortState) (i j : Nat) : LocalSortState :=
  { values := (state.values.set! i state.values[j]!).set! j state.values[i]!
    reversedCalls := (true, i, j) :: state.reversedCalls }

private def insertLeft (a : Nat) : Nat → LocalSortState → LocalSortState
  | 0, state => state
  | j + 1, state =>
    if j + 1 ≤ a then state else
      let (ordered, state) := lessAt state (j + 1) j
      if ordered then insertLeft a j (exchange state (j + 1) j) else state

/-- Go 1.27.0 src/sort/zsortinterface.go:10. Adjacent inversions move left. -/
def insertionSort (state : LocalSortState) (a b : Nat) : LocalSortState :=
  (List.range (b - (a + 1))).foldl (fun acc delta => insertLeft a (a + 1 + delta) acc) state

private def sift (hi first : Nat) : Nat → LocalSortState → Nat → LocalSortState
  | 0, state, _ => state
  | fuel + 1, state, root =>
    let child := 2 * root + 1
    if hi ≤ child then state else
      let (right, state) :=
        if child + 1 < hi then lessAt state (first + child) (first + child + 1)
        else (false, state)
      let greater := if right then child + 1 else child
      let (ordered, state) := lessAt state (first + root) (first + greater)
      if ordered then sift hi first fuel (exchange state (first + root) (first + greater)) greater
      else state

/-- Go 1.27.0 src/sort/zsortinterface.go:20. Each swap descends to a higher root index. -/
def siftDown (state : LocalSortState) (lo hi first : Nat) : LocalSortState :=
  sift hi first hi state lo

/-- Go 1.27.0 src/sort/zsortinterface.go:38. Build from right to left and pop through
index zero, including the final self-swap for a nonempty interval. -/
def heapSort (state : LocalSortState) (a b : Nat) : LocalSortState :=
  let hi := b - a
  let state := (List.range ((hi - 1) / 2 + 1)).reverse.foldl
    (fun acc root => siftDown acc root hi a) state
  (List.range hi).reverse.foldl
    (fun acc endIndex => siftDown (exchange acc a (a + endIndex)) 0 endIndex a) state

def insertionSortFunc := insertionSort
def siftDownFunc := siftDown
def heapSortFunc := heapSort

theorem insertionSort_empty (state : LocalSortState) (a : Nat) :
    insertionSort state a a = state := by simp [insertionSort, Nat.sub_eq_zero_of_le]

theorem siftDown_no_child (state : LocalSortState) (lo hi first : Nat) (h : hi ≤ 2 * lo + 1) :
    siftDown state lo hi first = state := by
  cases hi with
  | zero => rfl
  | succ n => simp [siftDown, sift, h]

theorem heapSort_empty (state : LocalSortState) (a : Nat) :
    heapSort state a a = state := by simp [heapSort, siftDown, sift]

theorem heapSort_single_calls (state : LocalSortState) :
    (heapSort state 0 1).reversedCalls = (true, 0, 0) :: state.reversedCalls := rfl

theorem insertionSortFunc_same (s : LocalSortState) (a b : Nat) :
    insertionSortFunc s a b = insertionSort s a b := rfl

theorem siftDownFunc_same (s : LocalSortState) (lo hi first : Nat) :
    siftDownFunc s lo hi first = siftDown s lo hi first := rfl

theorem heapSortFunc_same (s : LocalSortState) (a b : Nat) :
    heapSortFunc s a b = heapSort s a b := rfl

end GoLean.Sort
