import Init.Omega

namespace GoLean.Container.Heap

inductive Call where
  | less (child parent : Nat)
  | swap (parent child : Nat)

def up (child decisions : Nat) : List Call :=
  if child = 0 then []
  else
    let parent := (child - 1) / 2
    if decisions % 2 = 0 then [.less child parent]
    else .less child parent :: .swap parent child :: up parent (decisions / 2)
termination_by child
decreasing_by omega

theorem up_root (decisions : Nat) : up 0 decisions = [] := by
  rw [up]
  simp

theorem up_stop (child decisions : Nat)
    (positive : child ≠ 0) (stop : decisions % 2 = 0) :
    up child decisions = [.less child ((child - 1) / 2)] := by
  rw [up]
  simp [positive, stop]

theorem up_swap (child decisions : Nat)
    (positive : child ≠ 0) (proceed : decisions % 2 ≠ 0) :
    up child decisions = .less child ((child - 1) / 2) ::
      .swap ((child - 1) / 2) child :: up ((child - 1) / 2) (decisions / 2) := by
  rw [up]
  simp [positive, proceed]

end GoLean.Container.Heap
