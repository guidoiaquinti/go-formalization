import GoLean.Slices.PublicSort
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

private def floatKey (word : Nat) : Nat :=
  if word < 2 ^ 63 then word + 2 ^ 63 else 2 ^ 64 - 1 - word

private def floatLess (x y : BitVec 64) : Bool :=
  let a := x.toNat % 2 ^ 63
  let b := y.toNat % 2 ^ 63
  let nan := 0x7ff0000000000000
  if nan < a then decide (b ≤ nan)
  else if nan < b || (a = 0 && b = 0) then false
  else decide (floatKey x.toNat < floatKey y.toNat)

private def wire (o : PublicSortObservation (BitVec 64) Int) : String :=
  let status := match o.failure with
    | none => "ok"
    | some (.sequence .nilSequence) => "nilSequence"
    | some (.sequence (.producer code)) => s!"producer:{code}"
    | some (.sort (.callback code)) => s!"compare:{code}"
    | some (.sort .nilComparator) | some (.sort .nilCounter) => "nilCompare"
  let data := match o.values with
    | none => "nil"
    | some xs => "slice:" ++ String.intercalate "," (xs.map (toString ∘ BitVec.toNat))
  let calls := String.intercalate ";" (o.calls.map fun (x,y) => s!"{x.toNat},{y.toNat}")
  s!"{status}/{data}/{o.acknowledged}/{calls}"

def main : IO Unit := runTuple fun fields => do
  let width :: kind :: tag :: producer :: mode :: stop :: mask :: raw := fields |
    throw (IO.userError "expected native width, value kind, slice and iterator setup, comparator script and words")
  let words ← raw.mapM (wordInput 64)
  let input := if tag == 0 then none else some words
  let sequence := if producer == 0 then none else some (words,
    if producer == 2 then Except.error 7 else Except.ok ())
  let less := if kind == 1 then floatLess else fun x y : BitVec 64 => decide (x.toInt < y.toInt)
  let compare := if mode == 0 then none else some (fun (i : Nat) (x y : BitVec 64) =>
    if Int.ofNat i == stop then Except.error stop
    else if mode == 1 then Except.ok (if less x y then -1 else if less y x then 1 else 0)
    else Except.ok (if mask.toNat.testBit (i % 12) then -1 else 0))
  return String.intercalate "|" [
    wire (sliceSort width.toNat input less),
    wire (sliceSortFunc width.toNat input compare),
    wire (sliceSorted width.toNat sequence less),
    wire (sliceSortedFunc width.toNat sequence compare)]
