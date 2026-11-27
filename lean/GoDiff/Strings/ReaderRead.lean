import GoLean.Strings.ReaderRead
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

def main : IO Unit := runTuple fun fields => do
  let tag :: cursor :: previous :: dstLength :: raw := fields
    | throw (IO.userError "expected receiver, cursor, rune marker and destination length")
  let cursor ← wordInput 63 cursor
  let dst ← (raw.take dstLength.toNat).mapM (wordInput 8)
  let src ← (raw.drop dstLength.toNat).mapM (wordInput 8)
  let r := if tag = 0 then none else some (⟨src, cursor.toNat, previous⟩ : ReaderState)
  return match readerRead r dst with
  | none => "panic"
  | some (n, eof, bytes, s) =>
    String.intercalate " " (s!"{n} {eof} {s.cursor} {s.previousRune}" ::
      bytes.map (fun b => toString b.toNat))
