import GoLean.Encoding.Binary.CoderSkip
import GoDiff.Math.Bits.Tuple

open GoLean.Encoding.Binary

def main : IO Unit := runTuple fun row => do
  let width :: kind :: size :: offset :: raw := row |
    throw (IO.userError "expected word width, layout kind, size, offset and backing")
  let bytes ← raw.mapM (wordInput 8)
  let layout := match kind with
    | 0 => Layout.invalid
    | 1 => Layout.scalar 4
    | 2 => Layout.array size.toNat (.scalar 2)
    | 3 => Layout.slice size.toNat (.scalar 1)
    | 4 => Layout.structField (.scalar 1) (.structField (.scalar 8) .structEmpty)
    | 5 => Layout.structEmpty
    | 6 => Layout.array size.toNat .unsupported
    | _ => Layout.unsupported
  let wire := fun cells : List (BitVec 8) =>
    String.intercalate "," (cells.map (toString ∘ BitVec.toNat))
  let encoded := match encoderSkip width.toNat layout offset bytes with
    | none => s!"bounds:{offset}:{wire bytes}"
    | some (next, cells) => s!"ok:{next}:{wire cells}"
  return s!"{decoderSkip width.toNat layout offset}|{encoded}"
