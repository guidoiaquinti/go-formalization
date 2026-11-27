import GoLean.Bytes.ReaderRead
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def wire (xs : List (BitVec 8)) : String :=
  String.intercalate "," (xs.map (toString ∘ BitVec.toNat))

private def describe : Option ReaderState → String
  | none => "receiver-nil"
  | some s =>
    let marker := if s.bytes.isNone then "nil" else "slice"
    s!"{marker}:{s.cursor}:{s.previousRune}:{wire (s.bytes.getD [])}:{wire s.spare}"

private def render : Option ReaderCopyResult → String
  | none => "panic"
  | some r =>
    let status := match r.status with
      | .ok => "ok"
      | .eof => "EOF"
      | .negativeOffset => "bytes.Reader.ReadAt: negative offset"
    s!"{r.count}:{status}:{wire r.destination}:{describe r.state}"

def main : IO Unit := runTuple fun row => do
  let tag :: size :: spareSize :: cursor :: previous :: offset :: sharing :: dstSize :: raw := row |
    throw (IO.userError "expected reader and destination geometry")
  let data ← raw.mapM (wordInput 8)
  let capacity := size.toNat + spareSize.toNat
  let receiver : Option ReaderState :=
    if tag == 0 then none
    else some ⟨(if tag == 1 then none else some (data.take size.toNat)),
      (data.drop size.toNat).take spareSize.toNat, cursor.toNat, previous⟩
  let destination := if sharing < 0 then ReaderDestination.separate (data.drop capacity)
    else ReaderDestination.shared sharing.toNat dstSize.toNat
  return render (readerRead receiver destination) ++ "|" ++ render (readerReadAt receiver offset destination)
