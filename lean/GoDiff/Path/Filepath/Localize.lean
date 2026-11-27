import GoLean.Path.Filepath.Localize
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

private def readDevices : Nat → List Int → IO (List (List (BitVec 8)) × List Int)
  | 0, bytes => pure ([], bytes)
  | _ + 1, [] => throw (IO.userError "missing device length")
  | count + 1, length :: raw => do
    if length < 0 || length.toNat > raw.length then
      throw (IO.userError "invalid device length")
    let device ← (raw.take length.toNat).mapM (wordInput 8)
    let (devices, bytes) ← readDevices count (raw.drop length.toNat)
    return (device :: devices, bytes)

def main : IO Unit := runTuple fun fields => do
  let profile :: count :: fields := fields | throw (IO.userError "expected profile and device count")
  if count < 0 then throw (IO.userError "invalid device count")
  let profile := if profile == 1 then AbsProfile.windows
    else if profile == 2 then AbsProfile.plan9
    else AbsProfile.unix
  let (devices, bytes) ← readDevices count.toNat fields
  let path ← bytes.mapM (wordInput 8)
  match localizePath profile devices.contains path with
  | none => return "invalid path"
  | some result => return "ok:" ++ String.intercalate "," (result.map (toString ∘ BitVec.toNat))
