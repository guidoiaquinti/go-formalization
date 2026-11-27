import GoLean.Strings.ReaderSeek
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

def main : IO Unit := runTuple fun fields => do
  let tag :: cursor :: previous :: offset :: whence :: raw := fields
    | throw (IO.userError "expected receiver, seek parameters and backing bytes")
  let cursor ← wordInput 63 cursor
  let src ← raw.mapM (wordInput 8)
  let r := if tag = 0 then none else some (⟨src, cursor.toNat, previous⟩ : ReaderState)
  return match readerSeek r (BitVec.ofInt 64 offset) whence with
  | none => "panic"
  | some (pos, err, s) =>
    let status := match err with
      | none => "ok"
      | some .invalidWhence => "whence"
      | some .negativePosition => "negative"
    s!"{pos} {status} {s.cursor} {s.previousRune}"
