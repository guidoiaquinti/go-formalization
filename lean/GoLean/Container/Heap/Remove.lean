import GoLean.Container.Heap.Fix

namespace GoLean.Container.Heap

inductive RemoveCall where
  | len
  | initialSwap (index last : Nat)
  | repair (call : FixCall)
  | pop

/-- The caller's index is in range, so length - 1 is the former last index. -/
def remove (index length decisions value : Nat) : Nat × List RemoveCall :=
  let last := length - 1
  let calls :=
    if index = last then []
    else .initialSwap index last :: (fix index last decisions).tail.map .repair
  (value, .len :: calls ++ [.pop])

theorem remove_return (index length decisions value : Nat) :
    (remove index length decisions value).1 = value := rfl

theorem remove_last (length decisions value : Nat) :
    remove (length - 1) length decisions value = (value, [.len, .pop]) := by
  simp [remove]

theorem remove_repair (index length decisions value : Nat) (other : index ≠ length - 1) :
    (remove index length decisions value).2 =
      .len :: .initialSwap index (length - 1) ::
        ((fix index (length - 1) decisions).tail.map .repair ++ [.pop]) := by
  simp [remove, other]

end GoLean.Container.Heap
