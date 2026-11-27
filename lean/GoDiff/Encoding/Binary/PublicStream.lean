import GoLean.Encoding.Binary.PublicStream
import GoDiff.Math.Bits.Tuple

open GoLean.Encoding.Binary

private def kind : Int → FastKind
  | 0 => .bool | 1 => .int8 | 2 => .uint8 | 3 => .int16 | 4 => .uint16
  | 5 => .int32 | 6 => .uint32 | 7 => .int64 | 8 => .uint64
  | 9 => .float32 | 10 => .float64 | _ => .uint8

private def wire (values : List Nat) : String :=
  String.intercalate "," (values.map toString)

def main : IO Unit := runTuple fun row => do
  let tag :: form :: endian :: wordCount :: sourceLength :: chunk :: readTag :: writeTag :: raw := row |
    throw (IO.userError "expected type, form, order, size, input length, chunk and I/O scripts")
  let form := match form with
    | 0 => FastForm.scalar | 1 => .pointer | 2 => .slice | _ => .unsupported
  let words := (raw.take wordCount.toNat).map Int.toNat
  let source ← (raw.drop wordCount.toNat).mapM (wordInput 8)
  let source := source.take sourceLength.toNat
  let terminal := if readTag == 0 then "eof" else if readTag == 1 then "other" else "wrapped"
  let response := if writeTag == 0 then "ok" else if writeTag == 1 then "other" else "panic"
  let r := binaryRead (endian != 0) (tag == 11) (kind tag) form words source chunk.toNat terminal
  let w := binaryWrite (endian != 0) (tag == 11) (kind tag) form words response
  return s!"{r.status}:{wire r.state}:{wire r.requests}|{w.status}:{w.calls}:{w.aliases}:{wire (w.output.map BitVec.toNat)}"
