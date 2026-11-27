import GoLean.Strings.BuilderReset
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

def main : IO Unit := runTuple fun fields => do
  let receiver ← match fields with
    | [0] => pure none
    | 1 :: hasIdentity :: capacity :: input => do
      let bytes ← input.mapM (wordInput 8)
      let buffer := if capacity < 0 then none else some (bytes, capacity.toNat)
      pure (some { identitySet := hasIdentity != 0, buffer })
    | _ => throw (IO.userError "expected a nil marker or builder state")
  match builderReset receiver with
  | none => return "panic"
  | some result =>
    let (bytes, capacity) := result.buffer.getD ([], 0)
    return String.intercalate " " (toString capacity :: toString bytes.length ::
      bytes.map (fun byte => toString byte.toNat))
