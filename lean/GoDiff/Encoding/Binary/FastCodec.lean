import GoLean.Encoding.Binary.FastCodec
import GoDiff.Math.Bits.Tuple

open GoLean.Encoding.Binary

private def parseKind : Int → FastKind
  | 0 => .bool | 1 => .int8 | 2 => .uint8 | 3 => .int16 | 4 => .uint16
  | 5 => .int32 | 6 => .uint32 | 7 => .int64 | 8 => .uint64
  | 9 => .float32 | 10 => .float64 | _ => .unsupported

def main : IO Unit := runTuple fun row => do
  let tag :: shape :: endian :: count :: length :: raw := row |
    throw (IO.userError "expected type, shape, endian, count, length and cells")
  let form := match shape with
    | 0 => FastForm.scalar | 1 => FastForm.pointer | 2 => FastForm.slice | _ => FastForm.unsupported
  let big := endian != 0
  let n := count.toNat
  let m := length.toNat
  let values := (raw.take n).map Int.toNat
  let cells ← (raw.drop n).mapM (wordInput 8)
  let (read, decoded) := decodeFast big (parseKind tag) form (cells.take m) values
  let (written, encoded) := encodeFast big (parseKind tag) form values (cells.take m)
  let render := fun values : List Nat => String.intercalate "," (values.map toString)
  return s!"{read}:{render decoded}|{written}:{render ((encoded ++ cells.drop m).map BitVec.toNat)}"
