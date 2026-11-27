import Init.Data.BitVec.Basic

namespace GoLean.Encoding.Binary

private def wordBytes (bigEndian : Bool) (width word : Nat) : List (BitVec 8) :=
  (List.range width).map fun index =>
    let position := if bigEndian then width - 1 - index else index
    BitVec.ofNat 8 (word / 256 ^ position)

def putWord (bigEndian : Bool) (width word : Nat) (input : List (BitVec 8)) :
    Bool × List (BitVec 8) :=
  if input.length < width then (false, input)
  else (true, wordBytes bigEndian width word ++ input.drop width)

def bigPutUint16 (word : BitVec 16) := putWord true 2 word.toNat
def bigPutUint32 (word : BitVec 32) := putWord true 4 word.toNat
def bigPutUint64 (word : BitVec 64) := putWord true 8 word.toNat
def littlePutUint16 (word : BitVec 16) := putWord false 2 word.toNat
def littlePutUint32 (word : BitVec 32) := putWord false 4 word.toNat
def littlePutUint64 (word : BitVec 64) := putWord false 8 word.toNat

theorem putWord_short (bigEndian : Bool) (width word : Nat) (input : List (BitVec 8))
    (short : input.length < width) :
    putWord bigEndian width word input = (false, input) := by
  simp [putWord, short]

theorem putWord_suffix (bigEndian : Bool) (width word : Nat) (input : List (BitVec 8))
    (fits : width ≤ input.length) :
    (putWord bigEndian width word input).2.drop width = input.drop width := by
  simp [putWord, Nat.not_lt.mpr fits, wordBytes]

theorem bigPutUint16_short (word : BitVec 16) (input : List (BitVec 8))
    (short : input.length < 2) : bigPutUint16 word input = (false, input) :=
  putWord_short true 2 word.toNat input short

theorem bigPutUint16_suffix (word : BitVec 16) (input : List (BitVec 8))
    (fits : 2 ≤ input.length) : (bigPutUint16 word input).2.drop 2 = input.drop 2 :=
  putWord_suffix true 2 word.toNat input fits

theorem bigPutUint32_short (word : BitVec 32) (input : List (BitVec 8))
    (short : input.length < 4) : bigPutUint32 word input = (false, input) :=
  putWord_short true 4 word.toNat input short

theorem bigPutUint32_suffix (word : BitVec 32) (input : List (BitVec 8))
    (fits : 4 ≤ input.length) : (bigPutUint32 word input).2.drop 4 = input.drop 4 :=
  putWord_suffix true 4 word.toNat input fits

theorem bigPutUint64_short (word : BitVec 64) (input : List (BitVec 8))
    (short : input.length < 8) : bigPutUint64 word input = (false, input) :=
  putWord_short true 8 word.toNat input short

theorem bigPutUint64_suffix (word : BitVec 64) (input : List (BitVec 8))
    (fits : 8 ≤ input.length) : (bigPutUint64 word input).2.drop 8 = input.drop 8 :=
  putWord_suffix true 8 word.toNat input fits

theorem littlePutUint16_short (word : BitVec 16) (input : List (BitVec 8))
    (short : input.length < 2) : littlePutUint16 word input = (false, input) :=
  putWord_short false 2 word.toNat input short

theorem littlePutUint16_suffix (word : BitVec 16) (input : List (BitVec 8))
    (fits : 2 ≤ input.length) : (littlePutUint16 word input).2.drop 2 = input.drop 2 :=
  putWord_suffix false 2 word.toNat input fits

theorem littlePutUint32_short (word : BitVec 32) (input : List (BitVec 8))
    (short : input.length < 4) : littlePutUint32 word input = (false, input) :=
  putWord_short false 4 word.toNat input short

theorem littlePutUint32_suffix (word : BitVec 32) (input : List (BitVec 8))
    (fits : 4 ≤ input.length) : (littlePutUint32 word input).2.drop 4 = input.drop 4 :=
  putWord_suffix false 4 word.toNat input fits

theorem littlePutUint64_short (word : BitVec 64) (input : List (BitVec 8))
    (short : input.length < 8) : littlePutUint64 word input = (false, input) :=
  putWord_short false 8 word.toNat input short

theorem littlePutUint64_suffix (word : BitVec 64) (input : List (BitVec 8))
    (fits : 8 ≤ input.length) : (littlePutUint64 word input).2.drop 8 = input.drop 8 :=
  putWord_suffix false 8 word.toNat input fits

end GoLean.Encoding.Binary
