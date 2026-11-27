import GoLean.Strings.BuilderWriteByte
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

def main : IO Unit := runTuple fun fields => do
  let mode :: identity :: capacity :: byte :: raw := fields
    | throw (IO.userError "expected builder state and byte")
  let input ← wordInput 8 byte
  let contents ← raw.mapM (wordInput 8)
  let buffer := if capacity < 0 then none else some (contents, capacity.toNat)
  let receiver := if mode = 2 then none else some ({ identitySet := identity != 0, buffer } : BuilderState)
  return match builderWriteByte receiver (mode = 1) input with
  | none => "nil"
  | some (written, contents, identity) =>
    let status := if written then "ok" else "copied"
    String.intercalate " " (s!"{status} {identity}" ::
      contents.map (fun b => toString b.toNat))
