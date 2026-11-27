import GoLean.Go.Constant.Integer

namespace GoLean.Go.Constant

/-- Go 1.27.0 src/go/constant/value.go:185–186,247–248. -/
def int64ValString (word : BitVec 64) : String := toString word.toInt

def int64ValExactString (word : BitVec 64) : String := int64ValString word

def intValString (value : Option Int) : String := (value.map toString).getD "<nil>"

def intValExactString (value : Option Int) : String := intValString value

/-- Pinned signed-width predicates in value.go:965–975. -/
def is32bit (x : BitVec 64) : Bool := -2147483648 ≤ x.toInt && x.toInt ≤ 2147483647

def is63bit (x : BitVec 64) : Bool :=
  -4611686018427387904 ≤ x.toInt && x.toInt ≤ 4611686018427387903

/-- The fraction component threshold in value.go:347–355. -/
def smallInt (x : Nat) : Bool := (if x = 0 then 0 else x.log2 + 1) < 4096

theorem int64ValString_min : int64ValString 0x8000000000000000 = "-9223372036854775808" := by decide

theorem int64ValExactString_delegate (word : BitVec 64) : int64ValExactString word = int64ValString word := rfl

theorem intValString_nil : intValString none = "<nil>" := rfl

theorem intValExactString_delegate (value : Option Int) : intValExactString value = intValString value := rfl

theorem is32bit_edges : is32bit 0x7fffffff = true ∧ is32bit 0x80000000 = false := by decide

theorem is63bit_edges : is63bit 0x3fffffffffffffff = true ∧ is63bit 0x4000000000000000 = false := by decide

theorem smallInt_zero : smallInt 0 = true := by decide

end GoLean.Go.Constant
