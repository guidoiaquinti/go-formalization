import GoLean.Archive.Zip.EntryStat
import GoDiff.Math.Bits.Tuple

open GoLean.Archive.Zip

private def wire : StatResult → String
  | .header identity => s!"header:{identity}"
  | .directory identity => s!"directory:{identity}"
  | .duplicate raw => "duplicate:" ++ String.intercalate "," (raw.map (toString ·.toNat))
  | .nilAccess => "nil"

def main : IO Unit := runTuple fun fields => do
  let dup :: dir :: hasFile :: hasEntry :: words := fields |
    throw (IO.userError "expected duplicate, directory, file and entry flags and raw name")
  let name ← words.mapM fun b => return UInt8.ofNat (← wordInput 8 b).toNat
  let header := if hasFile == 1 then some 29 else none
  let e : StatEntry := ⟨17, name, dup == 1, dir == 1, header⟩
  return s!"{wire (listEntryStat e)}|{wire (openDirStat (if hasEntry == 1 then some e else none))}|{wire (checksumReaderStat header)}"
