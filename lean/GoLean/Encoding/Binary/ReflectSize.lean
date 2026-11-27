namespace GoLean.Encoding.Binary

inductive Layout where
  | scalar (width : Nat)
  | array (length : Nat) (element : Layout)
  | structEmpty
  | structField (field rest : Layout)
  | slice (length : Nat) (element : Layout)
  | pointer (null : Bool) (element : Layout)
  | unsupported
  | invalid

def typeSize : Layout → Option Nat
  | .scalar width => some width
  | .array length element => (typeSize element).map (length * ·)
  | .structEmpty => some 0
  | .structField field rest => do
    let n ← typeSize field
    let m ← typeSize rest
    pure (n + m)
  | _ => none

def sizeof : Layout → Except String (Option Nat)
  | .invalid => .error "nil type"
  | layout => .ok (typeSize layout)

def dataSize : Layout → Option Nat
  | .slice length element => (typeSize element).map (length * ·)
  | layout => typeSize layout

def binarySize : Layout → Option Nat
  | .pointer true _ => none
  | .pointer false element => dataSize element
  | layout => dataSize layout

theorem sizeof_unpadded (field rest : Layout) (f r : Nat)
    (fieldSize : typeSize field = some f) (restSize : typeSize rest = some r) :
    sizeof (.structField field rest) = .ok (some (f + r)) := by
  simp [sizeof, typeSize, fieldSize, restSize]

theorem sizeof_unsupported_element (n : Nat) (element : Layout)
    (bad : typeSize element = none) : sizeof (.array n element) = .ok none := by
  simp [sizeof, typeSize, bad]

theorem sizeof_nil : sizeof .invalid = .error "nil type" := rfl

theorem dataSize_slice (n : Nat) (element : Layout) :
    dataSize (.slice n element) = (typeSize element).map (n * ·) := rfl

theorem binarySize_nil (element : Layout) :
    binarySize (.pointer true element) = none := rfl

theorem binarySize_indirect (element : Layout) :
    binarySize (.pointer false element) = dataSize element := rfl

end GoLean.Encoding.Binary
