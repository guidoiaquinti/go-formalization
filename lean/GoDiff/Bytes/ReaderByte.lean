import GoLean.Bytes.ReaderByte
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def wire (xs : List (BitVec 8)) : String :=
  String.intercalate "," (xs.map (toString ∘ BitVec.toNat))

private def describe (state : ReaderState) : String :=
  let marker := if state.bytes.isNone then "nil" else "slice"
  s!"{marker}:{state.cursor}:{state.previousRune}:{wire (state.bytes.getD [])}:{wire state.spare}"

def main : IO Unit := runTuple fun row => do
  let tag :: size :: cursor :: previous :: raw := row | throw (IO.userError "expected reader state")
  let data ← raw.mapM (wordInput 8)
  let receiver : Option ReaderState :=
    if tag == 0 then none
    else some ⟨(if tag == 1 then none else some (data.take size.toNat)),
      data.drop size.toNat, cursor.toNat, previous⟩
  let read := match readerReadByte receiver with
    | none => "panic"
    | some (value, eof, state) => s!"{value.toNat}:{if eof then "EOF" else "ok"}:{describe state}"
  let unread := match readerUnreadByte receiver with
    | none => "panic"
    | some (err, state) => s!"{err.getD "ok"}:{describe state}"
  return s!"{read}|{unread}"
