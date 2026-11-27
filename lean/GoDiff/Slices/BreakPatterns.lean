import GoLean.Slices.BreakPatterns
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

private def wire (result : SwapRangeResult Int) : String :=
  let status := if result.boundsPanic then "panic" else "ok"
  match result.cells with
  | none => status ++ ":nil"
  | some values => status ++ ":slice:" ++ String.intercalate "," (values.map toString)

def main : IO Unit := runTuple fun fields => do
  let tag :: width :: a :: b :: length :: mode :: cells := fields |
    throw (IO.userError "expected nil marker, word width, range offsets, length, comparator and backing")
  if length < 0 || cells.length < length.toNat then
    throw (IO.userError "invalid backing geometry")
  let data := if tag == 0 then none else some cells
  let compare := if mode == 0 then none else some (fun x y : Int => x - y)
  return wire (breakPatternsOrdered width.toNat data length.toNat a b) ++ "|" ++
    wire (breakPatternsCmpFunc width.toNat data length.toNat a b compare)
