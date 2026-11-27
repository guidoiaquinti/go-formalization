import GoLean.Strings.AppendSliceWriterWrite
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

def main : IO Unit := runTuple fun fields => do
  let width :: capacity :: nilInput :: fields := fields
    | throw (IO.userError "expected receiver size, capacity and input nilness")
  let raw ← fields.mapM (wordInput 8)
  let receiver : Option AppendSliceState :=
    if width == -1 then none
    else some ⟨if width == -2 then none else some (raw.take width.toNat), capacity.toNat⟩
  let input := if nilInput == 1 then none else some (raw.drop width.toNat)
  let some result := appendSliceWriterWrite receiver input | return "panic"
  let bytes := result.bytes.getD []
  let nilness := if result.bytes.isNone then "nil" else "slice"
  let capacity := (result.capacity.map toString).getD "growth"
  return s!"{result.count}:{result.error}:{nilness}:{capacity}:" ++
    String.intercalate "," (bytes.map (fun byte => toString byte.toNat))
