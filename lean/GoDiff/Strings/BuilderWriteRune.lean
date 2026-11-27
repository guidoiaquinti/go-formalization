import GoLean.Strings.BuilderWriteRune
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

def main : IO Unit := runTuple fun fields => do
  let mode :: identity :: capacity :: rune :: raw := fields
    | throw (IO.userError "expected builder state and rune")
  let contents ← raw.mapM (wordInput 8)
  let buffer := if capacity < 0 then none else some (contents, capacity.toNat)
  let receiver := if mode = 2 then none else some ({ identitySet := identity != 0, buffer } : BuilderState)
  return match builderWriteRune receiver (mode = 1) rune with
  | none => "nil"
  | some result =>
    let status := (result.count.map fun n => s!"ok {n}").getD "copied"
    String.intercalate " " (s!"{status} {result.identitySet}" ::
      result.contents.map (fun b => toString b.toNat))
