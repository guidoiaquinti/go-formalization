import GoLean.Encoding.Binary.CoderWord
import GoDiff.Math.Bits.Tuple

open GoLean.Encoding.Binary

private def render (r : WordCoderResult) (writing : Bool) : String :=
  let value := if writing || r.status != "ok" then "0" else toString r.value
  let calls := (r.call.map toString).getD "-"
  s!"{r.status}:{value}:{r.offset}:{calls}:{String.intercalate "," (r.cells.map (toString ∘ BitVec.toNat))}"

def main : IO Unit := runTuple fun row => do
  let kind :: tag :: length :: offset :: arg :: raw := row |
    throw (IO.userError "need width, order, length, offset, word and backing")
  let cells ← raw.mapM (wordInput 8)
  let order := match tag with
    | 0 => WordOrder.little
    | 1 => WordOrder.big
    | 2 => WordOrder.scripted false
    | 3 => WordOrder.scripted true
    | _ => WordOrder.missing
  let s : ByteCoder := ⟨cells, length.toNat, offset⟩
  let results := match kind with
    | 2 => [decoderUint16 s order, decoderInt16 s order, encoderUint16 s order (BitVec.ofInt 16 arg),
         encoderInt16 s order (BitVec.ofInt 16 arg)]
    | 4 => [decoderUint32 s order, decoderInt32 s order, encoderUint32 s order (BitVec.ofInt 32 arg),
         encoderInt32 s order (BitVec.ofInt 32 arg)]
    | _ => [decoderUint64 s order, decoderInt64 s order, encoderUint64 s order (BitVec.ofInt 64 arg),
         encoderInt64 s order (BitVec.ofInt 64 arg)]
  return String.intercalate "|" ((results.zipIdx).map fun (result, i) => render result (i ≥ 2))
