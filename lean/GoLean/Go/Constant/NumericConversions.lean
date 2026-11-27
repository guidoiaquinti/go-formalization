namespace GoLean.Go.Constant

structure IntegerCell where
  identity : Nat
  value : Int

structure RationalCell where
  identity : Nat
  numerator : Int
  denominator : Nat

structure FloatCell where
  identity : Nat
  integerValue : Int
  precision : Nat

inductive ConstantRef where
  | absent
  | typedNil
  | identity (id : Nat)

structure ComplexCell where
  real : ConstantRef
  imaginary : Int

/-- Go 1.27.0 src/go/constant/value.go:273–283.
    Identities stand for the separately allocated numeric values. -/
def newNumericInt (identity : Nat) : IntegerCell := ⟨identity, 0⟩
def newNumericRat (identity : Nat) : RationalCell := ⟨identity, 0, 1⟩
def newNumericFloat (identity : Nat) : FloatCell := ⟨identity, 0, 512⟩

def i64toi (identity : Nat) (value : Int) : IntegerCell :=
  { newNumericInt identity with value }

def i64tor (identity : Nat) (value : Int) : RationalCell :=
  { newNumericRat identity with numerator := value }

def i64tof (identity : Nat) (value : Int) : FloatCell :=
  { newNumericFloat identity with integerValue := value }

def itor (identity : Nat) (source : IntegerCell) : RationalCell :=
  { newNumericRat identity with numerator := source.value }

def vtoc (source : ConstantRef) : ComplexCell := ⟨source, 0⟩

theorem newNumericInt_zero (id : Nat) : newNumericInt id = ⟨id, 0⟩ := rfl

theorem newNumericRat_zero (id : Nat) : newNumericRat id = ⟨id, 0, 1⟩ := rfl

theorem newNumericFloat_zero (id : Nat) : newNumericFloat id = ⟨id, 0, 512⟩ := rfl

theorem i64toi_exact (id : Nat) (value : Int) : i64toi id value = ⟨id, value⟩ := rfl

theorem i64tor_exact (id : Nat) (value : Int) : i64tor id value = ⟨id, value, 1⟩ := rfl

theorem i64tof_exact (id : Nat) (value : Int) : i64tof id value = ⟨id, value, 512⟩ := rfl

theorem itor_exact (id : Nat) (source : IntegerCell) :
    itor id source = ⟨id, source.value, 1⟩ := rfl

theorem vtoc_preserves (source : ConstantRef) : vtoc source = ⟨source, 0⟩ := rfl

end GoLean.Go.Constant
