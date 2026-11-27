import GoLean.Path.Filepath.Abs
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

private def render (bytes : List (BitVec 8)) : String :=
  String.intercalate "," (bytes.map (toString ∘ BitVec.toNat))

def main : IO Unit := runTuple fun fields => do
  let profile :: failure :: pathLength :: queryLength :: raw := fields |
    throw (IO.userError "expected profile, outcome, path length and query length")
  if pathLength < 0 || queryLength < 0 || pathLength.toNat + queryLength.toNat > raw.length then
    throw (IO.userError "invalid path or query length")
  let bytes ← raw.mapM (wordInput 8)
  let path := bytes.take pathLength.toNat
  let query := (bytes.drop pathLength.toNat).take queryLength.toNat
  let outcome := bytes.drop (pathLength.toNat + queryLength.toNat)
  let observed := if failure == 1 then Except.error outcome else Except.ok outcome
  let fullPath := fun arg => if arg == query then observed else
    Except.error ("unexpected FullPath input".toUTF8.data.toList.map UInt8.toBitVec)
  let profile := if profile == 1 then AbsProfile.windows else if profile == 2 then AbsProfile.plan9
    else AbsProfile.unix
  match absPath profile path observed fullPath with
  | .ok path => return "ok:" ++ render path
  | .error err => return "err:" ++ render err
