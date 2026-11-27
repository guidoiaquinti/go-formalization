import GoLean.Archive.Zip.FileMode
import GoDiff.Math.Bits.Tuple

open GoLean.Archive.Zip

def main : IO Unit := runTuple fun fields => do
  let [version, attrs, mode, flags, slash] := fields |
    throw (IO.userError "expected version, attributes, mode, flags and slash")
  let version ← wordInput 16 version
  let attrs ← wordInput 32 attrs
  let mode ← wordInput 32 mode
  let flags ← wordInput 16 flags
  let (newVersion, newAttrs) := fileHeaderSetMode version mode
  let descriptor : Nat := if fileHeaderHasDataDescriptor flags then 1 else 0
  return s!"{(fileHeaderMode version attrs (slash == 1)).toNat} {newVersion.toNat} {newAttrs.toNat} {descriptor} {(msdosModeToFileMode attrs).toNat} {(unixModeToFileMode attrs).toNat} {(fileModeToUnixMode mode).toNat}"
