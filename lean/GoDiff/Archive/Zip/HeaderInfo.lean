import GoLean.Archive.Zip.HeaderInfo
import GoDiff.Math.Bits.Tuple

open GoLean.Archive.Zip

def main : IO Unit := runTuple fun fields => do
  let creator :: attrs :: legacy :: wide :: present :: bytes := fields |
    throw (IO.userError "expected ZIP header words, presence, and name bytes")
  let creator ← wordInput 16 creator
  let attrs ← wordInput 32 attrs
  let legacy ← wordInput 32 legacy
  let wide ← wordInput 64 wide
  let id : Option Nat := if present == 1 then some 42 else none
  let info := fileHeaderFileInfo id
  let (returned, err) := headerInfoInfo info
  let view := s!"{info.header == id} {headerInfoSys info == id} {returned.header == id && err.isNone}"
  if id.isNone then return view
  let name := bytes.map Int.toNat
  let basename := String.intercalate "," ((headerInfoName name).map toString)
  return s!"{view} {basename} {headerInfoSize legacy wide} {(headerInfoMode creator attrs name).toNat} {(headerInfoType creator attrs name).toNat} {headerInfoIsDir creator attrs name}"
