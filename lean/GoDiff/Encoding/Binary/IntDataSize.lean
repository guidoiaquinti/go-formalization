import GoLean.Encoding.Binary.IntDataSize
import GoDiff.Math.Bits.Tuple

open GoLean.Encoding.Binary

private def kind : Int → FastKind
  | 0 => .bool | 1 => .int8 | 2 => .uint8 | 3 => .int16 | 4 => .uint16
  | 5 => .int32 | 6 => .uint32 | 7 => .int64 | 8 => .uint64
  | 9 => .float32 | 10 => .float64 | _ => .unsupported

def main : IO Unit := runTuple fun row => do
  let form :: token :: length :: nilTag :: raw := row |
    throw (IO.userError "expected form, kind, length, nil bit and backing")
  let backing ← raw.mapM (wordInput 8)
  let n := length.toNat
  let input := match form with
    | 0 => FastData.atom (kind token)
    | 1 | 2 => FastData.pointer (kind token)
    | 3 => if token == 2 then FastData.bytes ⟨backing, n, nilTag != 0⟩
        else FastData.elements (kind token) n
    | _ => FastData.unsupported
  let (size, returned) := intDataSize input
  return match returned with
    | none => s!"{size}:true:0:"
    | some view =>
      let cells := String.intercalate "," (view.backing.map (toString ∘ BitVec.toNat))
      s!"{size}:{view.nilSlice}:{view.backing.length}:{cells}"
