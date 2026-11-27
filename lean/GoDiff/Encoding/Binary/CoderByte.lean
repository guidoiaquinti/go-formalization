import GoLean.Encoding.Binary.CoderByte
import GoDiff.Math.Bits.Tuple

open GoLean.Encoding.Binary

private def wire (cells : List (BitVec 8)) : String :=
  String.intercalate "," (cells.map (toString ∘ BitVec.toNat))

def main : IO Unit := runTuple fun row => do
  let tag :: length :: offset :: arg :: raw := row |
    throw (IO.userError "need receiver, length, offset, byte and backing")
  if tag == 0 then return String.intercalate "|" (List.replicate 6 "nil")
  let cells ← raw.mapM (wordInput 8)
  let arg ← wordInput 8 arg
  let state : ByteCoder := ⟨cells, length.toNat, offset⟩
  let bounds := s!"bounds:{offset}:{wire cells}"
  let read {α : Type} [ToString α] (result : Option (α × Int)) : String :=
    result.elim bounds (fun (value, pos) => s!"{value}:{pos}:{wire cells}")
  let write (result : Option (List (BitVec 8) × Int)) : String :=
    result.elim bounds (fun (backing, pos) => s!"ok:{pos}:{wire backing}")
  return String.intercalate "|" [
    read (decoderBool state), read ((decoderUint8 state).map fun (b, p) => (b.toNat, p)),
    read (decoderInt8 state), write (encoderBool state (arg != 0)),
    write (encoderUint8 state arg), write (encoderInt8 state arg)]
