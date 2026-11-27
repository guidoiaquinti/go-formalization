import GoLean.Strings.ReaderUnreadRune
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

def main : IO Unit := runTuple fun fields => do
  let receiver ← match fields with
    | [0] => pure none
    | 1 :: cursor :: previous :: raw => do
      let cursor ← wordInput 63 cursor
      let bytes ← raw.mapM (wordInput 8)
      pure (some ⟨bytes, cursor.toNat, previous⟩)
    | _ => throw (IO.userError "expected nil marker or reader state")
  return match readerUnreadRune receiver with
    | none => "panic"
    | some (err, s) => s!"{s.cursor} {s.previousRune} {err.getD "ok"}"
