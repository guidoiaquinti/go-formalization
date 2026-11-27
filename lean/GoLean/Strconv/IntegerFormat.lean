import Init.Data.BitVec.Basic
import Init.Data.Int.DivMod

namespace GoLean.Strconv

private def digit (n : Nat) : Nat := if n < 10 then 48 + n else 87 + n

private def reverseDigits (base : Nat) : Nat → Nat → List Nat
  | 0, _ => []
  | fuel + 1, word =>
    if word < base then [digit word]
    else digit (word % base) :: reverseDigits base fuel (word / base)

private def radixDigits (base : Int) (word : Nat) : Option (List Nat) :=
  if base < 2 || base > 36 then none
  else some ((reverseDigits base.toNat 64 word).reverse)

/-- Go 1.27.0 internal/strconv/itoa.go formats the unsigned value in
lowercase, with decimal specialized and other bases delegated to formatBits. -/
def integerFormatUint (word : BitVec 64) (base : Int) : Option (List Nat) :=
  radixDigits base word.toNat

def integerFormatInt (word : BitVec 64) (base : Int) : Option (List Nat) :=
  (radixDigits base word.toInt.natAbs).map
    (fun digits => if word.toInt < 0 then 45 :: digits else digits)

def integerItoa (width : Nat) (word : Int) : Option (List Nat) :=
  integerFormatInt (BitVec.ofInt 64 (BitVec.ofInt width word).toInt) 10

structure FormatBuffer where
  cells : Option (List Nat)
  length : Nat

structure FormatAppend where
  output : Option (List Nat)
  original : Option (List Nat)
  reused : Bool
  panicked : Bool
  deriving DecidableEq

private def emit (buf : FormatBuffer) (neg : Bool) (digits : Option (List Nat)) : FormatAppend :=
  let old := buf.cells.getD []
  let staged := if neg then [45] else []
  let body := digits.getD []
  let out := old.take buf.length ++ staged ++ body
  let keepSign := buf.length + staged.length ≤ old.length
  let reused := keepSign && out.length ≤ old.length
  let kept := if reused then staged ++ body else if keepSign then staged else []
  let original := buf.cells.map (fun _ => old.take buf.length ++ kept ++ old.drop (buf.length + kept.length))
  match digits with
  | none => ⟨none, original, false, true⟩
  | some _ => ⟨some out, original, reused, false⟩

/-- AppendUint appends the complete digit sequence in a single operation. -/
def integerAppendUint (buf : FormatBuffer) (word : BitVec 64) (base : Int) : FormatAppend :=
  emit buf false (integerFormatUint word base)

/-- AppendInt writes a negative sign before delegating, even when radix validation fails. -/
def integerAppendInt (buf : FormatBuffer) (word : BitVec 64) (base : Int) : FormatAppend :=
  emit buf (word.toInt < 0) (radixDigits base word.toInt.natAbs)

def publicFormatUint := integerFormatUint
def publicFormatInt := integerFormatInt
def publicItoa := integerItoa
def publicAppendUint := integerAppendUint
def publicAppendInt := integerAppendInt

theorem integerFormatUint_zero : integerFormatUint 0 10 = some [48] := rfl
theorem integerFormatInt_minimum :
    integerFormatInt 0x8000000000000000 16 =
      some [45,56,48,48,48,48,48,48,48,48,48,48,48,48,48,48,48] := rfl
theorem integerItoa_negative : integerItoa 64 (-3) = some [45,51] := rfl
theorem integerAppendUint_badBase :
    integerAppendUint ⟨some [7,8], 1⟩ 3 1 = ⟨none, some [7,8], false, true⟩ := rfl
theorem integerAppendInt_badBase :
    integerAppendInt ⟨some [7,8], 1⟩ 0xffffffffffffffff 1 =
      ⟨none, some [7,45], false, true⟩ := rfl
theorem publicFormatUint_same (word : BitVec 64) (base : Int) :
    publicFormatUint word base = integerFormatUint word base := rfl
theorem publicFormatInt_same (word : BitVec 64) (base : Int) :
    publicFormatInt word base = integerFormatInt word base := rfl
theorem publicItoa_same (width : Nat) (word : Int) :
    publicItoa width word = integerItoa width word := rfl
theorem publicAppendUint_same (dst : FormatBuffer) (word : BitVec 64) (base : Int) :
    publicAppendUint dst word base = integerAppendUint dst word base := rfl
theorem publicAppendInt_same (dst : FormatBuffer) (word : BitVec 64) (base : Int) :
    publicAppendInt dst word base = integerAppendInt dst word base := rfl

end GoLean.Strconv
