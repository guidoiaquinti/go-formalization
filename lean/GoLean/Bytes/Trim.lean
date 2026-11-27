import GoLean.Strings.Trim
import GoLean.Bytes.IndexFuncPrivate

namespace GoLean.Bytes

private def leftOffset (input cutset : List (BitVec 8)) : Nat :=
  if cutset.isEmpty then 0 else
  let chars := cutset.map BitVec.toNat
  match indexFuncPrivate input (some fun r => GoLean.Strings.containsRune chars r) false with
  | some (i, _) => i.toNat
  | none => 0

/-- Go 1.27.0 src/bytes/bytes.go:987–1001. The shared right-then-left
    scanner determines retained bytes. Only the discarded left prefix
    reduces capacity, and an empty result is nil. -/
def trim (input cutset : List (BitVec 8)) (spare : Nat) :
    Option (List (BitVec 8) × Nat) :=
  let bytes := GoLean.Strings.trim input cutset
  if bytes.isEmpty then none
  else some (bytes, input.length + spare - leftOffset input cutset)

theorem trim_empty (cutset : List (BitVec 8)) (spare : Nat) :
    trim [] cutset spare = none := by
  simp [trim, GoLean.Strings.trim_emptyInput]

theorem trim_emptyCutset (head : BitVec 8) (tail : List (BitVec 8)) (spare : Nat) :
    trim (head :: tail) [] spare = some (head :: tail, (head :: tail).length + spare) := by
  simp [trim, GoLean.Strings.trim_emptyCutset, leftOffset]

theorem trim_retained (input cutset : List (BitVec 8)) (spare : Nat) :
    (trim input cutset spare).map Prod.fst =
      if (GoLean.Strings.trim input cutset).isEmpty then none
      else some (GoLean.Strings.trim input cutset) := by
  simp only [trim]
  split <;> simp_all

end GoLean.Bytes
