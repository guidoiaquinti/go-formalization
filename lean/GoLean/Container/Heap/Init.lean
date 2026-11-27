import GoLean.Container.Heap.Down

namespace GoLean.Container.Heap

inductive InitCall where
  | len
  | sift (call : Call)

private def isLess : Call → Bool
  | .less .. => true
  | .swap .. => false

def initSifts (parents length decisions : Nat) : List Call :=
  match parents with
  | 0 => []
  | index + 1 =>
    let calls := (down index length decisions).2
    calls ++ initSifts index length (decisions / 2 ^ calls.countP isLess)

def init (length decisions : Nat) : List InitCall :=
  .len :: (initSifts (length / 2) length decisions).map .sift

theorem init_len_first (length decisions : Nat) :
    (init length decisions).head? = some .len := rfl

theorem init_sifts (length decisions : Nat) :
    (init length decisions).tail = (initSifts (length / 2) length decisions).map .sift := rfl

theorem init_leaves (decisions : Nat) :
    init 0 decisions = [.len] ∧ init 1 decisions = [.len] := by
  exact ⟨rfl, rfl⟩

end GoLean.Container.Heap
