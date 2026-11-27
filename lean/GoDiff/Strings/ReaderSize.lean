import GoLean.Strings.ReaderSize
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let receiver ← match fields with
    | [0] => pure none
    | 1 :: cursor :: raw => do
      let cursor ← wordInput 63 cursor
      let bytes ← raw.mapM (wordInput 8)
      pure (some (bytes, cursor.toNat))
    | _ => throw (IO.userError "expected nil marker or cursor and bytes")
  return match GoLean.Strings.readerSize receiver with
    | none => "panic"
    | some size => toString size
