import GoLean.Container.Heap.Down

namespace GoLean.Container.Heap

inductive PopCall where
  | len
  | rootSwap (last : Int)
  | sift (call : Call)
  | pop

def pop (length decisions value : Nat) : Nat × List PopCall :=
  (value, .len :: .rootSwap ((length : Int) - 1) ::
    ((down 0 (length - 1) decisions).2.map .sift ++ [.pop]))

theorem pop_return (length decisions value : Nat) :
    (pop length decisions value).1 = value := rfl

theorem pop_dispatch (length decisions value : Nat) :
    (pop length decisions value).2.take 2 = [.len, .rootSwap ((length : Int) - 1)] := by
  rfl

theorem pop_sift (length decisions value : Nat) :
    (pop length decisions value).2.drop 2 =
      (down 0 (length - 1) decisions).2.map .sift ++ [.pop] := by
  rfl

theorem pop_empty_callbacks (decisions value : Nat) :
    pop 0 decisions value = (value, [.len, .rootSwap (-1), .pop]) := by
  simp [pop, down_leaf 0 0 decisions (by decide)]

end GoLean.Container.Heap
