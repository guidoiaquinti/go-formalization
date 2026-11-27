import GoLean.Path.Filepath.ReadDirNames
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

private def readNames : Nat → List Int → IO (List (List UInt8))
  | 0, [] => pure []
  | count + 1, size :: raw => do
    if size < 0 || size.toNat > raw.length then throw (IO.userError "invalid name length")
    let name ← (raw.take size.toNat).mapM fun value => do
      return UInt8.ofNat (← wordInput 8 value).toNat
    return name :: (← readNames count (raw.drop size.toNat))
  | _, _ => throw (IO.userError "invalid name list")

def main : IO Unit := runTuple fun fields => do
  let stage :: failure :: close :: count :: raw := fields |
    throw (IO.userError "expected stage, errors and name count")
  if count < 0 then throw (IO.userError "invalid name count")
  let names ← readNames count.toNat raw
  let observed :=
    if stage == 0 then DirNameRead.openError failure
    else DirNameRead.opened names (if failure == 0 then none else some failure)
      (if close == 0 then none else some close)
  let (trace, result) := readDirNames observed
  match result with
  | .error error => return s!"{trace.length}|nil|{error}"
  | .ok sorted =>
    let render := fun (name : List UInt8) => String.intercalate "," (name.map (toString ∘ UInt8.toNat))
    return s!"{trace.length}|{sorted.length}|{String.intercalate ";" (sorted.map render)}"
