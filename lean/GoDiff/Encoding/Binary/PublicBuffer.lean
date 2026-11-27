import GoLean.Encoding.Binary.PublicBuffer
import GoDiff.Math.Bits.Tuple

open GoLean.Encoding.Binary

private def kind : Int → FastKind
  | 0 => .bool | 1 => .int8 | 2 => .uint8 | 3 => .int16 | 4 => .uint16
  | 5 => .int32 | 6 => .uint32 | 7 => .int64 | 8 => .uint64
  | 9 => .float32 | 10 => .float64 | 11 => .uint8 | _ => .unsupported

private def byteWire (bs : List (BitVec 8)) : String :=
  String.intercalate "," (bs.map (toString ∘ BitVec.toNat))

def main : IO Unit := runTuple fun row => do
  let tag :: form :: big :: size :: length :: nilTag :: raw := row |
    throw (IO.userError "expected type, form, order, count, length, nil and data")
  let form := match form with
    | 0 => FastForm.scalar | 1 => FastForm.pointer | 2 => FastForm.slice | _ => FastForm.unsupported
  let words := (raw.take size.toNat).map Int.toNat
  let cells ← (raw.drop size.toNat).mapM (wordInput 8)
  let visible := cells.take length.toNat
  let d := binaryDecode (big != 0) (tag == 11) (kind tag) form visible words
  let e := binaryEncode (big != 0) (kind tag) form words visible
  let a := binaryAppend (big != 0) (kind tag) form words ⟨cells, length.toNat, nilTag != 0⟩
  let decoded := String.intercalate "," (d.state.map toString)
  return s!"{d.status}:{d.count}:{decoded}|{e.status}:{e.count}:{byteWire (e.state ++ cells.drop length.toNat)}|{a.status}:{a.allocated}:{a.nilOutput}:{byteWire a.output}:{byteWire a.backing}"
