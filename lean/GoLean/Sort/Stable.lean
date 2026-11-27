import GoLean.Sort.SymMerge
import GoLean.Sort.SortEntrypoints

namespace GoLean.Sort

private def insertStable (start : Nat) : Nat → MergeState → MergeState
  | 0, s => s
  | j + 1, s =>
    if j + 1 ≤ start then s
    else
      let (less, s) := mergeLess s (j + 1) j
      if less then insertStable start j (mergeSwap s (j + 1) j) else s

private def stableInsertion (s : MergeState) (a b : Nat) : MergeState :=
  (List.range (b - (a + 1))).foldl (fun s i => insertStable a (a + 1 + i) s) s

private def mergePass (s : MergeState) (n block : Nat) : MergeState :=
  let stride := 2 * block
  let full := n / stride
  let s := (List.range full).foldl (fun s k =>
    let a := k * stride
    symMerge s a (a + block) (a + stride)) s
  let tail := full * stride
  if tail + block < n then symMerge s tail (tail + block) n else s

private def stablePasses (n : Nat) : Nat → MergeState → Nat → MergeState
  | 0, s, _ => s
  | fuel + 1, s, block =>
    if block < n then stablePasses n fuel (mergePass s n block) (2 * block)
    else s

/-- Go 1.27.0 src/sort/zsortinterface.go:335. Insertion-sort blocks of 20,
then double the block width after each stable merge pass. -/
def stableSort (s : MergeState) (n : Int) : MergeState :=
  let count := n.toNat
  let full := count / 20
  let s := (List.range full).foldl (fun s k => stableInsertion s (20 * k) (20 * (k + 1))) s
  let s := stableInsertion s (20 * full) count
  stablePasses count count.log2 s 20

def stableSortFunc := stableSort

structure StableResult where
  status : String
  lenCalls : Nat
  state : MergeState

/-- Go 1.27.0 src/sort/sort.go:233. A nil interface panics during Len dispatch. -/
def stableInterface (s : MergeState) (reported : Option Int) : StableResult :=
  match reported with
  | none => ⟨"nil-interface", 0, s⟩
  | some count => ⟨"ok", 1, stableSort s count⟩

/-- Go 1.27.0 src/sort/slice.go:41. Swapper validates the input before sorting. -/
def stableSlice (s : MergeState) (shape : ReflectionSlice) (hasLess : Bool) : StableResult :=
  match shape with
  | .invalid => ⟨"invalid", 0, s⟩
  | .nonSlice => ⟨"non-slice", 0, s⟩
  | .valid =>
    if hasLess then ⟨"ok", 0, stableSortFunc s (Int.ofNat s.values.size)⟩
    else if s.values.size ≤ 1 then ⟨"ok", 0, s⟩ else ⟨"nil-less", 0, s⟩

theorem stableSort_zero (s : MergeState) : stableSort s 0 = s := rfl

theorem stableSort_single (s : MergeState) : stableSort s 1 = s := rfl

theorem stableSortFunc_same (s : MergeState) (n : Int) : stableSortFunc s n = stableSort s n := rfl

theorem stableInterface_delegates (s : MergeState) (n : Int) :
    stableInterface s (some n) = ⟨"ok", 1, stableSort s n⟩ := rfl

theorem stableSlice_invalid (s : MergeState) (hasLess : Bool) :
    stableSlice s .invalid hasLess = ⟨"invalid", 0, s⟩ := rfl

theorem stableSlice_short (s : MergeState) (h : s.values.size ≤ 1) :
    stableSlice s .valid false = ⟨"ok", 0, s⟩ := by simp [stableSlice, h]

end GoLean.Sort
