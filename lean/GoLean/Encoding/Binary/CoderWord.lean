import GoLean.Encoding.Binary.CoderByte
import GoLean.Encoding.Binary.EndianRead
import GoLean.Encoding.Binary.EndianPut

namespace GoLean.Encoding.Binary

inductive WordOrder
  | little
  | big
  | scripted (fail : Bool)
  | missing

structure WordCoderResult where
  status : String
  value : Int
  offset : Int
  cells : List (BitVec 8)
  call : Option Nat

private def runWord (state : ByteCoder) (order : WordOrder) (width : Nat)
    (input : Option Nat) : WordCoderResult :=
  let off := state.offset.toNat
  let failure := fun status => ⟨status, 0, state.offset, state.cells, none⟩
  -- The two-index slice uses capacity, unlike the single-byte index.
  if state.offset < 0 ∨ state.cells.length < off + width then failure "bounds"
  else
    match order with
    | .missing => failure "nil"
    | .scripted fail =>
      let cells := (state.cells.set off 165).set (state.cells.length - 1) 90
      let calls := some (state.cells.length - off)
      if fail then ⟨"callback", 0, state.offset, cells, calls⟩
      else ⟨"ok", Int.ofNat (256^width - 1), state.offset + width, cells, calls⟩
    | order =>
      let big := match order with | .big => true | _ => false
      let tail := state.cells.drop off
      let (value, cells) := match input with
        | none => ((readWord big width tail).getD 0, state.cells)
        | some word => (0, state.cells.take off ++ (putWord big width word tail).2)
      ⟨"ok", Int.ofNat value, state.offset + width, cells, none⟩

private def signedWord (width : Nat) (result : WordCoderResult) : WordCoderResult :=
  {result with value := (BitVec.ofInt (8 * width) result.value).toInt}

def decoderUint16 (s : ByteCoder) (o : WordOrder) := runWord s o 2 none
def decoderUint32 (s : ByteCoder) (o : WordOrder) := runWord s o 4 none
def decoderUint64 (s : ByteCoder) (o : WordOrder) := runWord s o 8 none
def decoderInt16 (s : ByteCoder) (o : WordOrder) := signedWord 2 (decoderUint16 s o)
def decoderInt32 (s : ByteCoder) (o : WordOrder) := signedWord 4 (decoderUint32 s o)
def decoderInt64 (s : ByteCoder) (o : WordOrder) := signedWord 8 (decoderUint64 s o)

def encoderUint16 (s : ByteCoder) (o : WordOrder) (x : BitVec 16) := runWord s o 2 (some x.toNat)
def encoderUint32 (s : ByteCoder) (o : WordOrder) (x : BitVec 32) := runWord s o 4 (some x.toNat)
def encoderUint64 (s : ByteCoder) (o : WordOrder) (x : BitVec 64) := runWord s o 8 (some x.toNat)
def encoderInt16 (s : ByteCoder) (o : WordOrder) (x : BitVec 16) := encoderUint16 s o x
def encoderInt32 (s : ByteCoder) (o : WordOrder) (x : BitVec 32) := encoderUint32 s o x
def encoderInt64 (s : ByteCoder) (o : WordOrder) (x : BitVec 64) := encoderUint64 s o x

theorem decoderUint16_bounds (s : ByteCoder) (o : WordOrder)
    (invalid : s.offset < 0 ∨ s.cells.length < s.offset.toNat + 2) :
    decoderUint16 s o = ⟨"bounds", 0, s.offset, s.cells, none⟩ := by
  simp [decoderUint16, runWord, invalid]

theorem decoderInt16_signed (s : ByteCoder) (o : WordOrder) :
    (decoderInt16 s o).value = (BitVec.ofInt 16 (decoderUint16 s o).value).toInt :=
  rfl

theorem encoderUint16_bounds (s : ByteCoder) (o : WordOrder) (x : BitVec 16)
    (invalid : s.offset < 0 ∨ s.cells.length < s.offset.toNat + 2) :
    encoderUint16 s o x = ⟨"bounds", 0, s.offset, s.cells, none⟩ := by
  simp [encoderUint16, runWord, invalid]

theorem encoderInt16_bits (s : ByteCoder) (o : WordOrder) (x : BitVec 16) :
    encoderInt16 s o x = encoderUint16 s o x :=
  rfl

theorem decoderUint32_bounds (s : ByteCoder) (o : WordOrder)
    (invalid : s.offset < 0 ∨ s.cells.length < s.offset.toNat + 4) :
    decoderUint32 s o = ⟨"bounds", 0, s.offset, s.cells, none⟩ := by
  simp [decoderUint32, runWord, invalid]

theorem decoderInt32_signed (s : ByteCoder) (o : WordOrder) :
    (decoderInt32 s o).value = (BitVec.ofInt 32 (decoderUint32 s o).value).toInt :=
  rfl

theorem encoderUint32_bounds (s : ByteCoder) (o : WordOrder) (x : BitVec 32)
    (invalid : s.offset < 0 ∨ s.cells.length < s.offset.toNat + 4) :
    encoderUint32 s o x = ⟨"bounds", 0, s.offset, s.cells, none⟩ := by
  simp [encoderUint32, runWord, invalid]

theorem encoderInt32_bits (s : ByteCoder) (o : WordOrder) (x : BitVec 32) :
    encoderInt32 s o x = encoderUint32 s o x :=
  rfl

theorem decoderUint64_bounds (s : ByteCoder) (o : WordOrder)
    (invalid : s.offset < 0 ∨ s.cells.length < s.offset.toNat + 8) :
    decoderUint64 s o = ⟨"bounds", 0, s.offset, s.cells, none⟩ := by
  simp [decoderUint64, runWord, invalid]

theorem decoderInt64_signed (s : ByteCoder) (o : WordOrder) :
    (decoderInt64 s o).value = (BitVec.ofInt 64 (decoderUint64 s o).value).toInt :=
  rfl

theorem encoderUint64_bounds (s : ByteCoder) (o : WordOrder) (x : BitVec 64)
    (invalid : s.offset < 0 ∨ s.cells.length < s.offset.toNat + 8) :
    encoderUint64 s o x = ⟨"bounds", 0, s.offset, s.cells, none⟩ := by
  simp [encoderUint64, runWord, invalid]

theorem encoderInt64_bits (s : ByteCoder) (o : WordOrder) (x : BitVec 64) :
    encoderInt64 s o x = encoderUint64 s o x :=
  rfl

end GoLean.Encoding.Binary
