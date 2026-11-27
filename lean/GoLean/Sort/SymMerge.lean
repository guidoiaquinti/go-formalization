import GoLean.Sort.RangeMoves

namespace GoLean.Sort

structure MergeState where
  values : Array (Int × Nat)
  reversedCalls : List (Bool × Nat × Nat)
  deriving Inhabited

def mergeLess (s : MergeState) (i j : Nat) : Bool × MergeState :=
  (s.values[i]!.1 < s.values[j]!.1,
   {s with reversedCalls := (false, i, j) :: s.reversedCalls})

def mergeSwap (s : MergeState) (i j : Nat) : MergeState :=
  { values := (s.values.set! i s.values[j]!).set! j s.values[i]!
    reversedCalls := (true, i, j) :: s.reversedCalls }

private def searchMerge (compare : MergeState → Nat → Bool × MergeState) :
    Nat → MergeState → Nat → Nat → Nat × MergeState
  | 0, s, lo, _ => (lo, s)
  | fuel + 1, s, lo, hi =>
    if lo < hi then
      let h := (lo + hi) / 2
      let (moveRight, s) := compare s h
      if moveRight then searchMerge compare fuel s (h + 1) hi
      else searchMerge compare fuel s lo h
    else (lo, s)

private def mergeMany : Nat → MergeState → Nat → Nat → Nat → MergeState
  | 0, s, _, _, _ => s
  | fuel + 1, s, a, m, b =>
    if m - a = 1 then
      let (i, s) := searchMerge (fun s h => mergeLess s h a) (b - m) s m b
      (List.range (i - 1 - a)).foldl
        (fun s k => mergeSwap s (a + k) (a + k + 1)) s
    else if b - m = 1 then
      let (i, s) := searchMerge (fun s h =>
        let (less, s) := mergeLess s m h
        (!less, s)) (m - a) s a m
      (List.range (m - i)).foldl
        (fun s k => mergeSwap s (m - k) (m - k - 1)) s
    else
      let mid := (a + b) / 2
      let n := mid + m
      let (lo, hi) := if mid < m then (n - b, mid) else (a, m)
      let (start, s) := searchMerge (fun s c =>
        let (less, s) := mergeLess s (n - 1 - c) c
        (!less, s)) (hi - lo) s lo hi
      let finish := n - start
      let s := if start < m && m < finish then
          (rotate start m finish).foldl (fun s (i, j) => mergeSwap s i j) s
        else s
      let s := if a < start && start < mid then mergeMany fuel s a start mid else s
      if mid < finish && finish < b then mergeMany fuel s mid finish b else s

/-- Go 1.27.0 src/sort/zsortinterface.go:378. Two nonempty sorted runs
have signed comparison keys and distinct, comparison-invisible identities. -/
def symMerge (s : MergeState) (a m b : Nat) : MergeState :=
  mergeMany (b - a) s a m b

def symMergeFunc := symMerge

theorem symMerge_two_ordered (s : MergeState) (h : ¬ s.values[1]!.1 < s.values[0]!.1) :
    (symMerge s 0 1 2).reversedCalls = (false, 1, 0) :: s.reversedCalls := by
  simp [symMerge, mergeMany, searchMerge, mergeLess, h]

theorem symMergeFunc_same (s : MergeState) (a m b : Nat) :
    symMergeFunc s a m b = symMerge s a m b := rfl

end GoLean.Sort
