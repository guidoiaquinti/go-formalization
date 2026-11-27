import GoLean.Slices.ReverseRange
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

private def wire (result : Option (List Int)) : String :=
  match result with
  | none => "nil"
  | some values => "slice:" ++ String.intercalate "," (values.map toString)

def main : IO Unit := runTuple fun fields => do
  let tag :: a :: b :: length :: mode :: cells := fields |
    throw (IO.userError "expected nil marker, range, length, comparator mode and backing")
  if a < 0 || b < a || length < b || cells.length < length.toNat then
    throw (IO.userError "invalid range or backing")
  let data := if tag == 0 then none else some cells
  let compare := if mode == 0 then none else some (fun x y : Int => x - y)
  return wire (reverseRangeOrdered data a.toNat b.toNat) ++ "|" ++
    wire (reverseRangeCmpFunc data a.toNat b.toNat compare)
