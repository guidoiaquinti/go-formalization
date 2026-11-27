namespace GoLean.Go.Constant

inductive ScalarValue where
  | unknown
  | boolean (value : Bool)
  | other (printed : String)

/-- Go 1.27.0 src/go/constant/value.go:381,384. -/
def makeUnknown : ScalarValue := .unknown
def makeBool (value : Bool) : ScalarValue := .boolean value

def unknownKind : Nat := 0
def boolKind (_ : Bool) : Nat := 1

def unknownString : String := "unknown"
def boolString (value : Bool) : String := if value then "true" else "false"

def unknownExactString : String := unknownString
def boolExactString (value : Bool) : String := boolString value

/-- Go 1.27.0 src/go/constant/value.go:478. Other values carry their
    separately rendered fmt argument. An error represents the panic text. -/
def boolVal : ScalarValue → Except String Bool
  | .boolean value => .ok value
  | .unknown => .ok false
  | .other representation => .error (representation ++ " not a Bool")

/-- Go 1.27.0 src/go/constant/kind_string.go:23. Known names are fixed slices of
    the generated name string. Outside those indices, render signed decimal. -/
def kindString (code : Int) : String :=
  let names := "UnknownBoolStringIntFloatComplex".toList
  let offsets := [0, 7, 11, 17, 20, 25, 32]
  if code < 0 || code ≥ 6 then "Kind(" ++ toString code ++ ")"
  else
    let start := offsets[code.toNat]?.getD 0
    let stop := offsets[code.toNat + 1]?.getD 0
    String.ofList ((names.drop start).take (stop - start))

theorem makeUnknown_default : boolVal makeUnknown = .ok false := rfl

theorem makeBool_roundTrip (b : Bool) : boolVal (makeBool b) = .ok b := rfl

theorem unknownKind_code : unknownKind = 0 := rfl

theorem boolKind_code (b : Bool) : boolKind b = 1 := rfl

theorem unknownString_text : unknownString = "unknown" := rfl

theorem boolString_text : boolString false = "false" ∧ boolString true = "true" := ⟨rfl, rfl⟩

theorem unknownExactString_short : unknownExactString = unknownString := rfl

theorem boolExactString_short (b : Bool) : boolExactString b = boolString b := rfl

theorem boolVal_wrongType (rendered : String) :
    boolVal (.other rendered) = .error (rendered ++ " not a Bool") := rfl

theorem kindString_negative (code : Int) (h : code < 0) :
    kindString code = "Kind(" ++ toString code ++ ")" := by
  simp [kindString, h]

theorem kindString_known :
    (List.range 6).map (fun code => kindString (Int.ofNat code)) =
      ["Unknown", "Bool", "String", "Int", "Float", "Complex"] := rfl

end GoLean.Go.Constant
