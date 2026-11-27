import GoLean.Slices.Insert
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

def main : IO Unit := runTuple fun fields => do
  let length :: idx :: nils :: cap :: alias :: count :: values := fields |
    throw (IO.userError "expected length, index, nilness, capacity and source information")
  let cells := values.take cap.toNat
  let source : InsertValues Int :=
    if alias < 0 then .independent (values.drop cap.toNat) else .alias alias.toNat count.toNat
  match sliceInsert ⟨cells, length.toNat, nils = 1⟩ idx source with
  | .error _ => return "bounds panic"
  | .ok out =>
    let header := s!"ok {out.view.length} {out.view.isNil} {out.allocated}"
    return String.intercalate " " (header :: out.view.cells.map toString)
