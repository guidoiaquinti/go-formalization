namespace GoLean.Strconv

structure NumError (Error : Type) where
  func : List UInt8
  num : List UInt8
  err : Error

/-- Go 1.27.0 src/strconv/number.go:251-262.
The error interface value is opaque. None for the receiver denotes a nil pointer. -/
def numErrorUnwrap (e : Option (NumError Error)) : Except Unit Error :=
  match e with
  | none => .error ()
  | some value => .ok value.err

theorem numErrorUnwrap_spec (e : Option (NumError Error)) (err : Error) :
    (numErrorUnwrap e = .ok err ↔ ∃ value, e = some value ∧ value.err = err) ∧
    (numErrorUnwrap e = .error () ↔ e = none) := by
  cases e <;> simp [numErrorUnwrap]


end GoLean.Strconv
