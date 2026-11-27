namespace GoLean.Encoding.Binary

structure ByteCoder where
  cells : List (BitVec 8)
  length : Nat
  offset : Int

def decoderUint8 (s : ByteCoder) : Option (BitVec 8 × Int) :=
  if s.offset < 0 ∨ s.length ≤ s.offset.toNat then none
  else (s.cells[s.offset.toNat]?).map fun byte => (byte, s.offset + 1)

def decoderBool (s : ByteCoder) : Option (Bool × Int) :=
  (decoderUint8 s).map fun (byte, offset) => (byte != 0, offset)

def decoderInt8 (s : ByteCoder) : Option (Int × Int) :=
  (decoderUint8 s).map fun (byte, offset) => (byte.toInt, offset)

def encoderUint8 (s : ByteCoder) (byte : BitVec 8) : Option (List (BitVec 8) × Int) :=
  if s.offset < 0 ∨ s.length ≤ s.offset.toNat then none
  else some (s.cells.set s.offset.toNat byte, s.offset + 1)

def encoderBool (s : ByteCoder) (value : Bool) : Option (List (BitVec 8) × Int) :=
  encoderUint8 s (if value then 1 else 0)

def encoderInt8 (s : ByteCoder) (byte : BitVec 8) : Option (List (BitVec 8) × Int) :=
  encoderUint8 s byte

theorem decoderUint8_in_bounds (s : ByteCoder) (byte : BitVec 8)
    (nonnegative : 0 ≤ s.offset) (bounds : s.offset.toNat < s.length)
    (cell : s.cells[s.offset.toNat]? = some byte) :
    decoderUint8 s = some (byte, s.offset + 1) := by
  simp [decoderUint8, Int.not_lt.mpr nonnegative, Nat.not_le.mpr bounds, cell]

theorem decoderBool_nonzero (s : ByteCoder) (byte : BitVec 8) (offset : Int)
    (read : decoderUint8 s = some (byte, offset)) :
    decoderBool s = some (byte != 0, offset) := by simp [decoderBool, read]

theorem decoderInt8_signed (s : ByteCoder) (byte : BitVec 8) (offset : Int)
    (read : decoderUint8 s = some (byte, offset)) :
    decoderInt8 s = some (byte.toInt, offset) := by simp [decoderInt8, read]

theorem encoderUint8_in_bounds (s : ByteCoder) (byte : BitVec 8)
    (nonnegative : 0 ≤ s.offset) (bounds : s.offset.toNat < s.length) :
    encoderUint8 s byte = some (s.cells.set s.offset.toNat byte, s.offset + 1) := by
  simp [encoderUint8, Int.not_lt.mpr nonnegative, Nat.not_le.mpr bounds]

theorem encoderBool_canonical (s : ByteCoder) (value : Bool) :
    encoderBool s value = encoderUint8 s (if value then 1 else 0) :=
  rfl

theorem encoderInt8_bits (s : ByteCoder) (byte : BitVec 8) :
    encoderInt8 s byte = encoderUint8 s byte :=
  rfl

end GoLean.Encoding.Binary
