import GoLean.Container.Heap.Up

namespace GoLean.Container.Heap

inductive PushCall where
  | push (value : Nat)
  | len
  | lift (call : Call)

def push (value length decisions : Nat) : List PushCall :=
  .push value :: .len :: (up (length - 1) decisions).map .lift

theorem push_dispatch (value length decisions : Nat) :
    (push value length decisions).take 2 = [.push value, .len] := by
  rfl

theorem push_sift (value length decisions : Nat) :
    (push value length decisions).drop 2 = (up (length - 1) decisions).map .lift := by
  rfl

theorem push_first (value decisions : Nat) :
    push value 1 decisions = [.push value, .len] := by
  simp [push, up_root]

end GoLean.Container.Heap
