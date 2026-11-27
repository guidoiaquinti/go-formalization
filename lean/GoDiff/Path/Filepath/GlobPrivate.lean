import GoLean.Path.Filepath.Glob
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

private def rawFields : Nat → List Int → IO (List (List (BitVec 8)))
  | 0, [] => pure []
  | count + 1, size :: raw => do
    if size < 0 || size.toNat > raw.length then throw (IO.userError "invalid field length")
    let value ← (raw.take size.toNat).mapM (wordInput 8)
    return value :: (← rawFields count (raw.drop size.toNat))
  | _, _ => throw (IO.userError "invalid field list")

def main : IO Unit := runTuple fun fields => do
  let windows :: state :: present :: collectedCount :: namesCount :: raw := fields |
    throw (IO.userError "expected profile, outcome and list counts")
  if collectedCount < 0 || namesCount < 0 then throw (IO.userError "invalid list count")
  let directory :: pattern :: tail ← rawFields (2 + collectedCount.toNat + namesCount.toNat) raw |
    throw (IO.userError "missing directory or pattern")
  let collected := if present == 0 then none else some (tail.take collectedCount.toNat)
  let directoryState :=
    if state == 0 then GlobDirectory.statFailed
    else if state == 1 then GlobDirectory.notDirectory
    else if state == 2 then GlobDirectory.openFailed
    else GlobDirectory.opened (tail.drop collectedCount.toNat)
  let (result, invalid) := globDirectory (windows == 1) directory pattern collected directoryState
  let status := if invalid then "bad" else "ok"
  match result with
  | none => return status ++ "|nil"
  | some paths =>
    let render := fun (path : List (BitVec 8)) =>
      String.intercalate "," (path.map (toString ∘ BitVec.toNat))
    return s!"{status}|{paths.length}:{String.intercalate ";" (paths.map render)}"
