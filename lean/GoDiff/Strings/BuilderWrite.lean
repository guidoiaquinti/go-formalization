import GoLean.Strings.BuilderWrite
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

def main : IO Unit := runTuple fun fields => do
  let mode :: identity :: capacity :: count :: nilSlice :: raw := fields
    | throw (IO.userError "expected builder state and input bytes")
  let contents ← (raw.take count.toNat).mapM (wordInput 8)
  let input ← (raw.drop count.toNat).mapM (wordInput 8)
  let buffer := if capacity < 0 then none else some (contents, capacity.toNat)
  let receiver := if mode = 2 then none else some ({ identitySet := identity != 0, buffer } : BuilderState)
  return match builderWrite receiver (mode = 1) (if nilSlice = 1 then none else some input) with
  | none => "nil"
  | some result =>
    let status := (result.count.map fun n => s!"ok {n}").getD "copied"
    String.intercalate " " (s!"{status} {result.identitySet}" ::
      result.contents.map (fun b => toString b.toNat))
