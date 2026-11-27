import GoLean.Encoding.Binary.VarintRead
import GoDiff.Math.Bits.Tuple

open GoLean.Encoding.Binary

private def wire : Option ReaderError → String
  | none => "ok"
  | some .eof => "eof"
  | some .unexpectedEOF => "unexpected"
  | some .overflow => "overflow"
  | some .other => "other"

def main : IO Unit := runTuple fun row => do
  let tag :: raw := row | throw (IO.userError "need error tag and reader bytes")
  let err := if tag == 0 then ReaderError.eof
    else if tag == 1 then ReaderError.unexpectedEOF else ReaderError.other
  let bytes ← raw.mapM (wordInput 8)
  let u := readUvarint bytes err
  let s := readVarint bytes err
  return s!"{u.value.toNat}:{wire u.error}:{u.calls}|{s.value}:{wire s.error}:{s.calls}"
