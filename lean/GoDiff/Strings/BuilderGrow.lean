import GoLean.Strings.BuilderGrow
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

def main : IO Unit := runTuple fun row => do
  let mode :: identity :: capacity :: request :: allocationRequest :: rounded :: raw := row |
    throw (IO.userError "expected builder state, request, allocation mapping and bytes")
  let bytes ← raw.mapM (wordInput 8)
  let buffer := if capacity < 0 then none else some (bytes, capacity.toNat)
  let receiver := if mode = 2 then none else some (⟨identity != 0, buffer⟩ : BuilderState)
  let allocate := fun input => if input = allocationRequest.toNat then rounded.toNat else input
  match builderGrow receiver (mode = 1) request allocate with
  | none => return "nil"
  | some (status, b) =>
    let tag := match status with
      | .ok => "ok"
      | .copied => "copied"
      | .negative => "negative"
    let (contents, cap) := b.buffer.getD ([], 0)
    return String.intercalate " " (s!"{tag} {b.identitySet} {cap}" ::
      contents.map (fun byte => toString byte.toNat))
