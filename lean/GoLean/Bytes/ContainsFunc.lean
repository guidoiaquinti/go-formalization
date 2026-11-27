import GoLean.Bytes.IndexFunc

namespace GoLean.Bytes

/-- Go 1.27.0 bytes/bytes.go:93–95. Testing the index preserves
    the callback trace and any nil-callback panic. -/
def containsFunc (input : List (BitVec 8)) (predicate : Option (Nat → Bool)) :
    Option (Bool × List Nat) :=
  (indexFunc input predicate).map fun (index, calls) => (decide (0 ≤ index), calls)

theorem containsFunc_empty (predicate : Option (Nat → Bool)) :
    containsFunc [] predicate = some (false, []) := rfl

theorem containsFunc_nil (lead : BitVec 8) (rest : List (BitVec 8)) :
    containsFunc (lead :: rest) none = none := rfl

theorem containsFunc_first (lead : BitVec 8) (rest : List (BitVec 8)) (f : Nat → Bool)
    (matched : f (GoLean.Unicode.Utf8.decodeNonempty lead.toNat (rest.map BitVec.toNat)).1 = true) :
    containsFunc (lead :: rest) (some f) =
      some (true, [(GoLean.Unicode.Utf8.decodeNonempty lead.toNat (rest.map BitVec.toNat)).1]) := by
  unfold containsFunc
  rw [indexFunc_first lead rest f matched]
  rfl

end GoLean.Bytes
