import GoLean.Strings.ReaderReadAt
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

def main : IO Unit := runTuple fun fields => do
  let tag :: offset :: dstLength :: raw := fields
    | throw (IO.userError "expected receiver marker, offset and destination length")
  let dst ← (raw.take dstLength.toNat).mapM (wordInput 8)
  let src ← (raw.drop dstLength.toNat).mapM (wordInput 8)
  let receiver := if tag = 0 then none else some (⟨src, 0, -1⟩ : ReaderState)
  return match readerReadAt receiver offset dst with
  | none => "panic"
  | some (n, error, dst) =>
    let status := match error with
      | none => "ok"
      | some .eof => "eof"
      | some .negativeOffset => "negative"
    String.intercalate " " (status :: toString n :: dst.map (fun b => toString b.toNat))
