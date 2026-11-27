import GoLean.Bytes.ReaderPosition
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def renderSeek : Option (Int × Option SeekError × ReaderState) → String
  | none => "nilReader"
  | some (pos, error, state) =>
    let status := match error with
      | none => "ok"
      | some .invalidWhence => "whence"
      | some .negativePosition => "negative"
    s!"{status}:{pos}:{state.cursor}:{state.previousRune}"

private def renderWrite : Option ReaderWriteResult → String
  | none => "nilReader"
  | some r =>
    let status := match r.status with
      | .ok => "ok"
      | .writerError => "writerError"
      | .shortWrite => "shortWrite"
      | .writerPanic => "writerPanic"
      | .invalidCount => "invalidCount"
      | .nilWriter => "nilWriter"
    let trace := match r.call with
      | none => "none"
      | some call =>
        s!"{call.capacity}:{String.intercalate "," (call.bytes.map (toString ∘ BitVec.toNat))}"
    s!"{status}:{r.count}:{r.cursor}:{r.previousRune}:{trace}"

def main : IO Unit := runTuple fun row => do
  let tag :: length :: cursor :: previous :: offset :: whence :: script :: count :: raw := row |
    throw (IO.userError "expected reader, seek and callback inputs")
  let data ← raw.mapM (wordInput 8)
  let receiver : Option ReaderState :=
    if tag == 0 then none
    else some ⟨(if tag == 1 then none else some (data.take length.toNat)),
      data.drop length.toNat, cursor.toNat, previous⟩
  let response : ReaderWriterResponse :=
    if script == 2 then .panicked else if script == 3 then .missing else .returned count (script == 1)
  return renderSeek (readerSeek receiver (BitVec.ofInt 64 offset) whence) ++ "|" ++
    renderWrite (readerWriteTo receiver response)
