import GoLean.Strings.BuilderGrowPrivate
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

def main : IO Unit := runTuple fun row => do
  let mode :: capacity :: request :: allocationRequest :: rounded :: raw := row |
    throw (IO.userError "expected builder state, request, allocation mapping and bytes")
  let bytes ← raw.mapM (wordInput 8)
  let buffer := if capacity < 0 then none else some (bytes, capacity.toNat)
  let receiver := if mode = 3 then none else some (⟨mode != 0, buffer⟩ : BuilderState)
  let allocate := fun input => if input = allocationRequest.toNat then rounded.toNat else input
  match builderGrowPrivate receiver request allocate with
  | none => return "nil"
  | some (status, b) =>
    let tag := match status with
      | .ok => "ok"
      | .allocationLength => "allocation"
      | .sliceBounds => "bounds"
    let identity := if !b.identitySet then "unset" else if mode = 1 then "self" else "foreign"
    let cap : Int := (b.buffer.map fun x => Int.ofNat x.2).getD (-1)
    let contents := b.buffer.map Prod.fst |>.getD []
    return s!"{tag}:{identity}:{cap}:" ++ String.intercalate "," (contents.map fun byte => toString byte.toNat)
