import GoLean.Strings.ReaderReset
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

def main : IO Unit := runTuple fun fields => do
  let (receiver, raw) ← match fields with
    | 0 :: raw => pure (none, raw)
    | 1 :: cursor :: previous :: oldLength :: raw => do
      let old ← (raw.take oldLength.toNat).mapM (wordInput 8)
      pure (some ⟨old, cursor.toNat, previous⟩, raw.drop oldLength.toNat)
    | _ => throw (IO.userError "expected receiver state followed by replacement bytes")
  let bytes ← raw.mapM (wordInput 8)
  return match readerReset receiver bytes with
    | none => "panic"
    | some state => String.intercalate " " (toString state.cursor ::
        toString state.previousRune :: state.contents.map (fun b => toString b.toNat))
