import Init.Data.BitVec.Basic

namespace GoLean.Encoding.Binary

def readWord (bigEndian : Bool) (width : Nat) (input : List (BitVec 8)) : Option Nat :=
  if input.length < width then none
  else
    let bytes := input.take width
    let ordered := if bigEndian then bytes else bytes.reverse
    some (ordered.foldl (fun word byte => 256 * word + byte.toNat) 0)

def bigUint16 := readWord true 2
def bigUint32 := readWord true 4
def bigUint64 := readWord true 8
def littleUint16 := readWord false 2
def littleUint32 := readWord false 4
def littleUint64 := readWord false 8

theorem readWord_short (bigEndian : Bool) (width : Nat) (input : List (BitVec 8))
    (short : input.length < width) : readWord bigEndian width input = none := by
  simp [readWord, short]

theorem readWord_prefix (bigEndian : Bool) (width : Nat) (input : List (BitVec 8)) :
    readWord bigEndian width input = readWord bigEndian width (input.take width) := by
  rw [readWord, readWord, List.length_take, List.take_take, Nat.min_self]
  by_cases short : input.length < width
  · have lower : min width input.length < width := by omega
    simp [short, lower]
  · have lower : ¬ min width input.length < width := by omega
    simp [short, lower]

theorem bigUint16_short (input : List (BitVec 8)) (short : input.length < 2) :
    bigUint16 input = none :=
  readWord_short true 2 input short

theorem bigUint16_prefix (input : List (BitVec 8)) :
    bigUint16 input = bigUint16 (input.take 2) :=
  readWord_prefix true 2 input

theorem bigUint32_short (input : List (BitVec 8)) (short : input.length < 4) :
    bigUint32 input = none :=
  readWord_short true 4 input short

theorem bigUint32_prefix (input : List (BitVec 8)) :
    bigUint32 input = bigUint32 (input.take 4) :=
  readWord_prefix true 4 input

theorem bigUint64_short (input : List (BitVec 8)) (short : input.length < 8) :
    bigUint64 input = none :=
  readWord_short true 8 input short

theorem bigUint64_prefix (input : List (BitVec 8)) :
    bigUint64 input = bigUint64 (input.take 8) :=
  readWord_prefix true 8 input

theorem littleUint16_short (input : List (BitVec 8)) (short : input.length < 2) :
    littleUint16 input = none :=
  readWord_short false 2 input short

theorem littleUint16_prefix (input : List (BitVec 8)) :
    littleUint16 input = littleUint16 (input.take 2) :=
  readWord_prefix false 2 input

theorem littleUint32_short (input : List (BitVec 8)) (short : input.length < 4) :
    littleUint32 input = none :=
  readWord_short false 4 input short

theorem littleUint32_prefix (input : List (BitVec 8)) :
    littleUint32 input = littleUint32 (input.take 4) :=
  readWord_prefix false 4 input

theorem littleUint64_short (input : List (BitVec 8)) (short : input.length < 8) :
    littleUint64 input = none :=
  readWord_short false 8 input short

theorem littleUint64_prefix (input : List (BitVec 8)) :
    littleUint64 input = littleUint64 (input.take 8) :=
  readWord_prefix false 8 input

end GoLean.Encoding.Binary
