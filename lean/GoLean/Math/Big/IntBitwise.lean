import Init.Data.Int.Bitwise
import Init.Data.Nat.Bitwise.Lemmas

namespace GoLean.Math.Big

private def natAndNot (a b : Nat) : Nat := a ^^^ (a &&& b)

def intAnd : Int → Int → Int
  | .ofNat a, .ofNat b => .ofNat (a &&& b)
  | .ofNat a, .negSucc b => .ofNat (natAndNot a b)
  | .negSucc a, .ofNat b => .ofNat (natAndNot b a)
  | .negSucc a, .negSucc b => .negSucc (a ||| b)

def intNot (a : Int) : Int := ~~~a

def intAndNot (a b : Int) : Int := intAnd a (intNot b)

def intOr : Int → Int → Int
  | .ofNat a, .ofNat b => .ofNat (a ||| b)
  | .ofNat a, .negSucc b => .negSucc (natAndNot b a)
  | .negSucc a, .ofNat b => .negSucc (natAndNot a b)
  | .negSucc a, .negSucc b => .negSucc (a &&& b)

def intXor : Int → Int → Int
  | .ofNat a, .ofNat b | .negSucc a, .negSucc b => .ofNat (a ^^^ b)
  | .ofNat a, .negSucc b | .negSucc a, .ofNat b => .negSucc (a ^^^ b)

def intLsh (x : Int) (bits : Nat) : Int := x <<< bits

def intRsh (x : Int) (bits : Nat) : Int := x >>> bits

inductive BitFault where
  | negativeIndex
  | invalidBit
  deriving DecidableEq

def intBit (x index : Int) : Except BitFault Nat :=
  if index < 0 then .error .negativeIndex else .ok ((intRsh x index.toNat).natAbs % 2)

def intSetBit (x index : Int) (bit : Nat) : Except BitFault Int :=
  if index < 0 then .error .negativeIndex
  else if bit > 1 then .error .invalidBit
  else
    let mask : Int := 1 <<< index.toNat
    .ok (if bit == 0 then intAndNot x mask else intOr x mask)

private def countZeros : Nat → Nat → Nat
  | 0, _ => 0
  | fuel + 1, n =>
    if n == 0 || n % 2 == 1 then 0 else 1 + countZeros fuel (n / 2)

def intTrailingZeroBits (x : Int) : Nat := countZeros (x.natAbs.log2 + 1) x.natAbs

theorem intAnd_idempotent (x : Int) : intAnd x x = x := by
  cases x <;> simp [intAnd]

theorem intNot_involution (x : Int) : intNot (intNot x) = x := by
  cases x <;> rfl

theorem intAndNot_self (x : Int) : intAndNot x x = 0 := by
  cases x <;> change Int.ofNat (natAndNot _ _) = 0 <;> simp [natAndNot]

theorem intOr_idempotent (x : Int) : intOr x x = x := by
  cases x <;> simp [intOr]

theorem intXor_self (x : Int) : intXor x x = 0 := by
  cases x <;> simp [intXor]

theorem intLsh_zero (x : Int) : intLsh x 0 = x := Int.shiftLeft_zero x

theorem intRsh_zero (x : Int) : intRsh x 0 = x := Int.shiftRight_zero x

theorem intBit_negative (x index : Int) (h : index < 0) :
    intBit x index = .error .negativeIndex := by simp [intBit, h]

theorem intSetBit_negative (x index : Int) (b : Nat) (h : index < 0) :
    intSetBit x index b = .error .negativeIndex := by simp [intSetBit, h]

theorem intSetBit_badValue (x index : Int) (b : Nat) (h : ¬ index < 0) (invalid : b > 1) :
    intSetBit x index b = .error .invalidBit := by simp [intSetBit, h, invalid]

theorem intTrailingZeroBits_sign (x : Int) :
    intTrailingZeroBits (-x) = intTrailingZeroBits x := by simp [intTrailingZeroBits]

theorem intTrailingZeroBits_zero : intTrailingZeroBits 0 = 0 := rfl

end GoLean.Math.Big
