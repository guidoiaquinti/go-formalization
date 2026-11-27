import GoLean.Archive.Zip.FileListEntry
import GoDiff.Math.Bits.Tuple

open GoLean.Archive.Zip

private def bytesWire (s : EntryName) : String :=
  String.intercalate "," (s.map (toString ∘ BitVec.toNat))

def main : IO Unit := runTuple fun row => do
  let n :: hasFile :: seconds :: nanos :: cells := row |
    throw (IO.userError "expected name lengths and entry metadata")
  let words ← cells.mapM (wordInput 8)
  let name := words.take n.toNat
  let other := words.drop n.toNat
  let t := if hasFile == 1 then some ⟨seconds, nanos.toNat⟩ else none
  let e : ListEntry := ⟨42, name, t⟩
  let part := splitEntryName name
  let mod := listEntryModTime e
  let (identity, failure) := listEntryInfo e
  return String.intercalate "|" [bytesWire part.directory, bytesWire part.element,
    toString part.trailingSlash, toString (fileEntryCompare name other),
    bytesWire (listEntryName e), toString (listEntrySize e), toString (listEntryMode e).toNat,
    toString (listEntryType e).toNat, toString (listEntryIsDir e), toString (listEntrySys e).isNone,
    toString mod.seconds, toString mod.nanoseconds, toString (identity == e.identity && failure.isNone),
    bytesWire (listEntryString e)]
