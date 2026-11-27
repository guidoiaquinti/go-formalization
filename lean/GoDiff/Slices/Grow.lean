import GoLean.Slices.Grow
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

def main : IO Unit := runTuple fun fields => do
  let n :: length :: isNil :: cells := fields |
    throw (IO.userError "expected growth, length, nilness and capacity cells")
  match sliceGrow 0 ⟨cells, length.toNat, isNil = 1⟩ n with
  | .error msg => return "panic " ++ msg
  | .ok result =>
    let view := result.view
    let header := s!"ok {view.length} {view.isNil} {result.allocated}"
    return String.intercalate " " (header :: view.cells.map toString)
