import GoLean.Encoding.Binary.ReflectSize

namespace GoLean.Encoding.Binary

private def skipSize (layout : Layout) : Int :=
  ((dataSize layout).map Int.ofNat).getD (-1)

def decoderSkip (width : Nat) (layout : Layout) (offset : Int) : Int :=
  (BitVec.ofInt width (offset + skipSize layout)).toInt

def encoderSkip (width : Nat) (layout : Layout) (offset : Int) (cells : List (BitVec 8)) :
    Option (Int × List (BitVec 8)) :=
  let next := decoderSkip width layout offset
  if offset < 0 ∨ next < offset ∨ cells.length < next.toNat then none
  else some (next, cells.take offset.toNat ++
    List.replicate (next.toNat - offset.toNat) 0 ++ cells.drop next.toNat)

theorem decoderSkip_unsupported (width : Nat) (offset : Int) :
    decoderSkip width .unsupported offset = (BitVec.ofInt width (offset - 1)).toInt :=
  rfl

theorem encoderSkip_bounds (width : Nat) (layout : Layout) (offset : Int)
    (cells : List (BitVec 8))
    (invalid : offset < 0 ∨ decoderSkip width layout offset < offset ∨
      cells.length < (decoderSkip width layout offset).toNat) :
    encoderSkip width layout offset cells = none := by
  simp only [encoderSkip, invalid, ite_true]

end GoLean.Encoding.Binary
