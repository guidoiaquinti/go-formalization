import GoLean.Strings.ReaderWriteTo
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

def main : IO Unit := runTuple fun fields => do
  let live :: cursor :: previous :: stringWriter :: script :: count :: bytes := fields
    | throw (IO.userError "expected reader state and writer script")
  let cursor ← wordInput 63 cursor
  let contents ← bytes.mapM (wordInput 8)
  let r := if live = 0 then none else some (⟨contents, cursor.toNat, previous⟩ : ReaderState)
  let response := match script with
    | 2 => WriterResult.panicked
    | 3 => WriterResult.missing
    | _ => WriterResult.returned count (script = 1)
  return match readerWriteTo r (stringWriter != 0) response with
  | none => "nilReader"
  | some s =>
    let status := match s.status with
      | .ok => "ok"
      | .writerError => "writerError"
      | .shortWrite => "shortWrite"
      | .writerPanic => "writerPanic"
      | .invalidCount => "invalidCount"
      | .nilWriter => "nilWriter"
    let trace := match s.call with
      | none => ["none"]
      | some (kind, payload) => (if kind then "string" else "bytes") ::
          payload.map (fun b => toString b.toNat)
    String.intercalate " " (s!"{status} {s.count} {s.cursor} {s.previousRune}" :: trace)
