import GoLean.Slices.RotateRange
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

private def wire (result : Except Unit (Option (List Int))) : String :=
  match result with
  | .error _ => "outside"
  | .ok none => "nil"
  | .ok (some values) => "slice:" ++ String.intercalate "," (values.map toString)

def main : IO Unit := runTuple fun fields => do
  let tag :: a :: middle :: b :: length :: mode :: cells := fields |
    throw (IO.userError "expected nil marker, range, length, comparator and backing")
  if min a (min middle b) < 0 || length < 0 || cells.length < length.toNat then
    throw (IO.userError "invalid backing geometry")
  let data := if tag == 0 then none else some cells
  let compare := if mode == 0 then none else some (fun x y : Int => x - y)
  return wire (rotateOrdered data length.toNat a.toNat middle.toNat b.toNat) ++ "|" ++
    wire (rotateCmpFunc data length.toNat a.toNat middle.toNat b.toNat compare)
