import Init.Data.BitVec.Basic

namespace GoLean.Math.Big

structure IntWords (width : Nat) where
  negative : Bool
  words : List (BitVec width)

private def wordParts (width : Nat) (value : BitVec 64) : List (BitVec width) :=
  if value = 0 then []
  else if width = 32 && value.toNat / 2^32 > 0 then
    [BitVec.ofNat width value.toNat, BitVec.ofNat width (value.toNat / 2^32)]
  else [BitVec.ofNat width value.toNat]

/-- Go 1.27.0 src/math/big/int.go:58–66. Set the sign and compact magnitude. -/
def setInt64 (width : Nat) (present : Bool) (value : BitVec 64) : Option (IntWords width) :=
  if !present then none else
    let negative := value.toInt < 0
    let magnitude := if negative then -value else value
    some ⟨negative, wordParts width magnitude⟩

/-- Go 1.27.0 int.go:70–74. Clear the sign and store the unsigned value. -/
def setUint64 (width : Nat) (present : Bool) (value : BitVec 64) : Option (IntWords width) :=
  if present then some ⟨false, wordParts width value⟩ else none

/-- Go 1.27.0 int.go:77–91. Fresh storage has the same normalized word contents. -/
def newInt (width : Nat) (value : BitVec 64) : IntWords width :=
  let negative := value.toInt < 0
  let u := if negative then -value else value
  ⟨negative, wordParts width u⟩

/-- Go 1.27.0 int.go:44–55. The empty magnitude always has sign zero. -/
def intSign (x : Option (IntWords width)) : Option Int :=
  x.map fun n => if n.words.isEmpty then 0 else if n.negative then -1 else 1

private def low64 (words : List (BitVec width)) : BitVec 64 :=
  let low := BitVec.ofNat 64 ((words.headD 0).toNat)
  if width = 32 then
    low ||| (BitVec.ofNat 64 ((words[1]?.getD 0).toNat) <<< 32)
  else low

/-- Go 1.27.0 int.go:512–518. Capture the pinned low signed word, including
    inexact values that the Go API leaves undefined. -/
def intToInt64 (x : Option (IntWords width)) : Option Int :=
  x.map fun n =>
    let low := low64 n.words
    (if n.negative then -low else low).toInt

/-- Go 1.27.0 int.go:522–524. Return the low word of the magnitude. -/
def intToUint64 (x : Option (IntWords width)) : Option Nat :=
  x.map fun n => (low64 n.words).toNat

/-- Go 1.27.0 int.go:527–533. MinInt64 is accepted only for a negative value. -/
def intIsInt64 (x : Option (IntWords width)) : Option Bool :=
  x.map fun n =>
    if n.words.length ≤ 64 / width then
      let word := low64 n.words
      word.toInt ≥ 0 || (n.negative && word == -word)
    else false

/-- Go 1.27.0 int.go:536–538. Check native word count and sign. -/
def intIsUint64 (x : Option (IntWords width)) : Option Bool :=
  x.map fun n => !n.negative && n.words.length ≤ 64 / width

/-- Go 1.27.0 int.go:631–636. Count full lower words and top-word bits. -/
def intBitLen (x : Option (IntWords width)) : Option Nat :=
  x.map fun n =>
    match n.words.getLast? with
    | none => 0
    | some top => (n.words.length - 1) * width +
        (if top = 0 then 0 else top.toNat.log2 + 1)

theorem setInt64_min : setInt64 32 true 0x8000000000000000 =
    some ⟨true, [0, 0x80000000]⟩ := rfl
theorem setUint64_max : setUint64 32 true 0xffffffffffffffff =
    some ⟨false, [0xffffffff, 0xffffffff]⟩ := rfl
theorem newInt_zero (width : Nat) : newInt width 0 = ⟨false, []⟩ := rfl
theorem intSign_empty (negative : Bool) :
    intSign (some (⟨negative, []⟩ : IntWords 64)) = some 0 := rfl
theorem intToInt64_negative : intToInt64 (some (⟨true, [1]⟩ : IntWords 64)) = some (-1) := rfl
theorem intToUint64_negative : intToUint64 (some (⟨true, [1]⟩ : IntWords 64)) = some 1 := rfl
theorem intIsInt64_min :
    intIsInt64 (some (⟨true, [0x8000000000000000]⟩ : IntWords 64)) = some true ∧
    intIsInt64 (some (⟨false, [0x8000000000000000]⟩ : IntWords 64)) = some false := ⟨rfl, rfl⟩
theorem intIsUint64_negative : intIsUint64 (some (⟨true, [1]⟩ : IntWords 64)) = some false := rfl
theorem intBitLen_min : intBitLen (some (⟨true, [0, 0x80000000]⟩ : IntWords 32)) = some 64 := rfl

end GoLean.Math.Big
