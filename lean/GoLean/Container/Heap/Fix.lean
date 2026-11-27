import GoLean.Container.Heap.Down

namespace GoLean.Container.Heap

inductive FixCall where
  | len
  | sift (call : Call)

private def comparisonCount (calls : List Call) : Nat :=
  calls.countP fun call => match call with
    | .less .. => true
    | .swap .. => false

def fix (index length decisions : Nat) : List FixCall :=
  let (moved, downCalls) := down index length decisions
  let upCalls :=
    if moved then []
    else up index (decisions / 2 ^ comparisonCount downCalls)
  .len :: (downCalls ++ upCalls).map .sift

theorem fix_moved (index length decisions : Nat) (calls : List Call)
    (moved : down index length decisions = (true, calls)) :
    fix index length decisions = .len :: calls.map .sift := by
  simp [fix, moved]

theorem fix_still (index length decisions : Nat) (calls : List Call)
    (still : down index length decisions = (false, calls)) :
    fix index length decisions =
      .len :: (calls ++ up index (decisions / 2 ^ comparisonCount calls)).map .sift := by
  simp [fix, still]

theorem fix_singleton (decisions : Nat) : fix 0 1 decisions = [.len] := by
  simp [fix, down_leaf 0 1 decisions (by decide), up_root]

end GoLean.Container.Heap
