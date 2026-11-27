import GoLean.Path.Filepath.UnixAbs
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

private def render (bytes : List (BitVec 8)) : String :=
  String.intercalate "," (bytes.map (toString ∘ BitVec.toNat))

def main : IO Unit := runTuple fun fields => do
  let profile :: failure :: size :: raw := fields |
    throw (IO.userError "expected profile, working directory status and path length")
  if size < 0 || size.toNat > raw.length then
    throw (IO.userError "invalid path length")
  let bytes ← raw.mapM (wordInput 8)
  let tail := bytes.drop size.toNat
  let profile := if profile == 1 then AbsProfile.windows else if profile == 2 then AbsProfile.plan9
    else AbsProfile.unix
  let cwd := if failure == 1 then Except.error tail else Except.ok tail
  match unixAbsPath profile (bytes.take size.toNat) cwd with
  | .ok path => return "ok:" ++ render path
  | .error err => return "err:" ++ render err
