import GoLean.Slices.Replace
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

def main : IO Unit := runTuple fun fields => do
  let width :: length :: lo :: hi :: nils :: cap :: alias :: count :: values := fields |
    throw (IO.userError "expected width, length, bounds, nilness, capacity and source information")
  let cells := values.take cap.toNat
  let source : InsertValues Int :=
    if alias < 0 then .independent (values.drop cap.toNat) else .alias alias.toNat count.toNat
  match sliceReplace 0 width.toNat ⟨cells, length.toNat, nils = 1⟩ lo hi source with
  | .error _ => return "bounds panic"
  | .ok out =>
    let header := s!"ok {out.view.length} {out.view.isNil} {out.allocated}"
    return String.intercalate " " (header :: out.view.cells.map toString)
