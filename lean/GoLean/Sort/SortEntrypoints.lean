import GoLean.Sort.Pdqsort

namespace GoLean.Sort

structure SortResult where
  status : String
  lenCalls : Nat
  state : LocalSortState

private def limitFor (n : Nat) : Nat :=
  if n = 0 then 0 else n.log2 + 1

/-- Go 1.27.0 src/sort/sort.go:50. Interface dispatch calls Len once before
    testing the early return. A missing interface panics during dispatch. -/
def sortInterface (nativeWidth : Nat) (s : LocalSortState) (reported : Option Int) : SortResult :=
  match reported with
  | none => ⟨"nil-interface", 0, s⟩
  | some n =>
    if n ≤ 1 then ⟨"ok", 1, s⟩
    else ⟨"ok", 1, pdqsort nativeWidth s 0 n.toNat (limitFor n.toNat)⟩

/-- Go 1.27.0 src/sort/sort.go:130. The receiver reports its actual length. -/
def intSliceSort (nativeWidth : Nat) (s : LocalSortState) : SortResult :=
  sortInterface nativeWidth s (some (Int.ofNat s.values.size))

inductive ReflectionSlice where
  | valid
  | invalid
  | nonSlice

/-- Go 1.27.0 src/sort/slice.go:24. Swapper validates the kind before any
    comparison. A valid slice of length at most one causes no Less call. -/
def sortSlice (nativeWidth : Nat) (s : LocalSortState) (shape : ReflectionSlice)
    (hasLess : Bool) : SortResult :=
  match shape with
  | .invalid => ⟨"invalid", 0, s⟩
  | .nonSlice => ⟨"non-slice", 0, s⟩
  | .valid =>
    let n := s.values.size
    if hasLess then ⟨"ok", 0, pdqsortFunc nativeWidth s 0 n (limitFor n)⟩
    else if n ≤ 1 then ⟨"ok", 0, s⟩
    else ⟨"nil-less", 0, s⟩

theorem sortInterface_short (w : Nat) (s : LocalSortState) (n : Int) (h : n ≤ 1) :
    sortInterface w s (some n) = ⟨"ok", 1, s⟩ := by simp [sortInterface, h]

theorem intSliceSort_delegates (w : Nat) (s : LocalSortState) :
    intSliceSort w s = sortInterface w s (some (Int.ofNat s.values.size)) := rfl

theorem sortSlice_invalid (w : Nat) (s : LocalSortState) (less : Bool) :
    sortSlice w s .invalid less = ⟨"invalid", 0, s⟩ := rfl

theorem sortSlice_nil_less_short (w : Nat) (s : LocalSortState) (h : s.values.size ≤ 1) :
    sortSlice w s .valid false = ⟨"ok", 0, s⟩ := by
  simp [sortSlice, h]

end GoLean.Sort
