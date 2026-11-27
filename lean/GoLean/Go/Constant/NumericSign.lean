import Init.Data.BitVec.Basic

namespace GoLean.Go.Constant

inductive SignValue where
  | int64 (value : Int)
  | bigInt (value : Int)
  | rational (numerator : Int)
  | float (signedMagnitude : Int)
  | complex (re im : SignValue)
  | unknown
  | other (printed : String)

/-- Go 1.27.0 src/go/constant/value.go:106–110. -/
def int64ValKind (_ : Int) : Nat := 3
def intValKind (_ : Int) : Nat := 3
def ratValKind (_ : Int) : Nat := 4
def floatValKind (_ : Int) : Nat := 4
def complexValKind (_ _ : SignValue) : Nat := 5

private def signedOrder (n : Int) : Int :=
  if n < 0 then -1 else if n = 0 then 0 else 1

/-- Go 1.27.0 src/go/constant/value.go:687–709. The operands of
    the complex bitwise OR are always -1, 0 or 1, so two bits suffice. -/
def constantSign : SignValue → Except String Int
  | .int64 n | .bigInt n | .rational n | .float n => .ok (signedOrder n)
  | .complex real imaginary => do
    let r ← constantSign real
    let i ← constantSign imaginary
    return (BitVec.ofInt 2 r ||| BitVec.ofInt 2 i).toInt
  | .unknown => .ok 1
  | .other text => .error (text ++ " not numeric")

theorem constantSign_unknown : constantSign .unknown = .ok 1 := rfl

theorem constantSign_mixed (a b : Int) (negative : a < 0) (positive : 0 < b) :
    constantSign (.complex (.int64 a) (.rational b)) = .ok (-1) := by
  have bNonnegative : ¬ b < 0 := Int.not_lt.mpr (Int.le_of_lt positive)
  have bNonzero : b ≠ 0 := Int.ne_of_gt positive
  simp [constantSign, signedOrder, negative, bNonnegative, bNonzero]
  rfl

theorem constantSign_wrongType (rendered : String) :
    constantSign (.other rendered) = .error (rendered ++ " not numeric") := rfl

theorem int64ValKind_code (x : Int) : int64ValKind x = 3 := rfl
theorem intValKind_code (x : Int) : intValKind x = 3 := rfl
theorem ratValKind_code (x : Int) : ratValKind x = 4 := rfl
theorem floatValKind_code (x : Int) : floatValKind x = 4 := rfl
theorem complexValKind_code (x y : SignValue) : complexValKind x y = 5 := rfl

end GoLean.Go.Constant
