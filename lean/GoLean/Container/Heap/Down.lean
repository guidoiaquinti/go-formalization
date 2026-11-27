import Init.Omega
import GoLean.Container.Heap.Up

namespace GoLean.Container.Heap

def down (parent length decisions : Nat) : Bool × List Call :=
  let left := 2 * parent + 1
  if left ≥ length then (false, [])
  else
    let right := left + 1
    let hasRight := right < length
    let child := if hasRight ∧ decisions % 2 ≠ 0 then right else left
    let first := if hasRight then [Call.less right left] else []
    let rest := if hasRight then decisions / 2 else decisions
    if rest % 2 = 0 then (false, first ++ [.less child parent])
    else
      (true, first ++ [.less child parent, .swap parent child] ++
        (down child length (rest / 2)).2)
termination_by length - parent
decreasing_by
  dsimp
  split <;> omega

theorem down_leaf (parent length decisions : Nat) (leaf : length ≤ 2 * parent + 1) :
    down parent length decisions = (false, []) := by
  rw [down]
  simp [leaf]

theorem down_left_stop (parent length decisions : Nat)
    (left : 2 * parent + 1 < length) (noRight : length ≤ 2 * parent + 1 + 1)
    (stop : decisions % 2 = 0) :
    down parent length decisions = (false, [.less (2 * parent + 1) parent]) := by
  rw [down]
  simp [Nat.not_le.mpr left, Nat.not_lt.mpr noRight, stop]

theorem down_right_swap (parent length decisions : Nat)
    (right : 2 * parent + 1 + 1 < length)
    (preferRight : decisions % 2 ≠ 0) (swap : (decisions / 2) % 2 ≠ 0) :
    down parent length decisions =
      (true, [.less (2 * parent + 1 + 1) (2 * parent + 1),
        .less (2 * parent + 1 + 1) parent, .swap parent (2 * parent + 1 + 1)] ++
        (down (2 * parent + 1 + 1) length ((decisions / 2) / 2)).2) := by
  have left : ¬ length ≤ 2 * parent + 1 := by omega
  rw [down]
  simp [left, right, preferRight, swap]

end GoLean.Container.Heap
