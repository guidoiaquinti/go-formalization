import Init.Data.BitVec.Basic

namespace GoLean.Go.Constant

inductive IntegerValue where
  | small (word : BitVec 64)
  | large (integer : Int)
  | unknown
  | other
  deriving DecidableEq

def makeInt64 (word : BitVec 64) : IntegerValue := .small word

def makeInt (x : Int) : IntegerValue :=
  if -(2^63) ≤ x ∧ x < 2^63 then .small (BitVec.ofInt 64 x) else .large x

def makeUint64 (word : BitVec 64) : IntegerValue :=
  if word.toNat < 2^63 then .small word else .large (Int.ofNat word.toNat)

def int64Val : IntegerValue → Option (Int × Bool)
  | .small word => some (word.toInt, true)
  | .large x => some ((BitVec.ofInt 64 x).toInt, false)
  | .unknown => some (0, false)
  | .other => none

def uint64Val : IntegerValue → Option (Nat × Bool)
  | .small word => some (word.toNat, 0 ≤ word.toInt)
  | .large x => some (x.natAbs % 2^64, 0 ≤ x ∧ x < 2^64)
  | .unknown => some (0, false)
  | .other => none

private def bitSize (n : Nat) : Nat := if n = 0 then 0 else n.log2 + 1

def bitLen : IntegerValue → Option Nat
  | .small word => some (bitSize word.toInt.natAbs)
  | .large x => some (bitSize x.natAbs)
  | .unknown => some 0
  | .other => none

private def little (n : Nat) : List UInt8 :=
  (List.range ((bitSize n + 7)/8)).map (fun i => UInt8.ofNat (n / 256^i))

def integerBytes : IntegerValue → Option (List UInt8)
  | .small word => some (little word.toInt.natAbs)
  | .large x => some (little x.natAbs)
  | .unknown | .other => none

def makeFromBytes (raw : List UInt8) : IntegerValue :=
  makeInt (Int.ofNat (raw.foldr (fun byte next => byte.toNat + 256 * next) 0))

theorem makeInt64_kind (x : BitVec 64) : makeInt64 x = .small x := rfl
theorem makeInt_overflow (x : Int) (h : 2^63 ≤ x) : makeInt x = .large x := by
  simp [makeInt] <;> omega
theorem makeUint64_max : makeUint64 0xffffffffffffffff = .large 18446744073709551615 := rfl
theorem int64Val_unknown : int64Val .unknown = some (0, false) := rfl
theorem uint64Val_unknown : uint64Val .unknown = some (0, false) := rfl
theorem bitLen_unknown : bitLen .unknown = some 0 := rfl
theorem integerBytes_unknown : integerBytes .unknown = none := rfl
theorem makeFromBytes_empty : makeFromBytes [] = .small 0 := rfl

end GoLean.Go.Constant
