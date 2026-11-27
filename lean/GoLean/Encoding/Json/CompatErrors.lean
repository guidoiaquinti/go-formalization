namespace GoLean.Encoding.Json

abbrev ErrorBytes := List Nat

private def ascii (s : String) : ErrorBytes := s.toList.map Char.toNat

structure ErrorType where
  pointer : Bool
  text : ErrorBytes

inductive ErrorOutcome (α : Type) where
  | ok (value : α)
  | nilPanic
  | nestedPanic
  deriving DecidableEq

inductive NestedError where
  | absent
  | text (identity : Nat) (value : ErrorBytes)
  | panic (identity : Nat)

private def nestedID : NestedError → Nat
  | .absent => 0
  | .text id _ | .panic id => id

structure CompatibilityError where
  typ : Option ErrorType
  nested : NestedError
  value : ErrorBytes
  structName : ErrorBytes
  field : ErrorBytes
  source : ErrorBytes

def syntaxError (e : Option CompatibilityError) : ErrorOutcome ErrorBytes :=
  match e with
  | none => .nilPanic
  | some e => .ok e.value

def unsupportedValueError (e : Option CompatibilityError) : ErrorOutcome ErrorBytes :=
  match e with
  | none => .nilPanic
  | some e => .ok (ascii "json: unsupported value: " ++ e.value)

def unsupportedTypeError (e : Option CompatibilityError) : ErrorOutcome ErrorBytes :=
  match e.bind CompatibilityError.typ with
  | none => .nilPanic
  | some typ => .ok (ascii "json: unsupported type: " ++ typ.text)

def invalidUnmarshalError (e : Option CompatibilityError) : ErrorOutcome ErrorBytes :=
  match e with
  | none => .nilPanic
  | some e =>
    match e.typ with
    | none => .ok (ascii "json: Unmarshal(nil)")
    | some typ => .ok (ascii (if typ.pointer then "json: Unmarshal(nil " else "json: Unmarshal(non-pointer ")
      ++ typ.text ++ [41])

def marshalerError (e : Option CompatibilityError) : ErrorOutcome ErrorBytes :=
  match e with
  | some ⟨some typ, .text _ message, _, _, _, source⟩ =>
    let method := if source.isEmpty then ascii "MarshalJSON" else source
    .ok (ascii "json: error calling " ++ method ++ ascii " for type " ++ typ.text ++ [58, 32] ++ message)
  | some ⟨some _, .panic _, _, _, _, _⟩ => .nestedPanic
  | _ => .nilPanic

private def unwrap (e : Option CompatibilityError) : ErrorOutcome Nat :=
  match e with
  | none => .nilPanic
  | some e => .ok (nestedID e.nested)

def marshalerErrorUnwrap := unwrap

def unmarshalTypeErrorUnwrap := unwrap

private def lastComponentNumeric (field : ErrorBytes) : Bool :=
  let suffix := field.reverse.takeWhile (· != 46)
  !suffix.isEmpty && suffix.all fun d => 48 ≤ d && d ≤ 57

def unmarshalTypeError (e : Option CompatibilityError) : ErrorOutcome ErrorBytes :=
  match e with
  | none => .nilPanic
  | some e =>
    match e.typ with
    | none => .nilPanic
    | some typ =>
      let target :=
        if e.structName.isEmpty && e.field.isEmpty then ascii "Go value"
        else
          let kind := if lastComponentNumeric e.field then [] else ascii "Go struct field "
          kind ++ e.structName ++ [46] ++ e.field
      let description := ascii "json: cannot unmarshal " ++ e.value ++ ascii " into " ++ target ++
        ascii " of type " ++ typ.text
      match e.nested with
      | .absent => .ok description
      | .text _ message => .ok (description ++ [58, 32] ++ message)
      | .panic _ => .nestedPanic

theorem syntaxError_value (e : CompatibilityError) : syntaxError (some e) = .ok e.value := rfl

theorem unsupportedValueError_raw (e : CompatibilityError) :
    unsupportedValueError (some e) = .ok (ascii "json: unsupported value: " ++ e.value) := rfl

theorem unsupportedTypeError_missing (e : CompatibilityError) (h : e.typ = none) :
    unsupportedTypeError (some e) = .nilPanic := by simp [unsupportedTypeError, h]

theorem invalidUnmarshalError_missing (e : CompatibilityError) (h : e.typ = none) :
    invalidUnmarshalError (some e) = .ok (ascii "json: Unmarshal(nil)") := by
  simp [invalidUnmarshalError, h]

theorem marshalerError_default (typ : ErrorType) (n : Nat) (message : ErrorBytes) :
    marshalerError (some ⟨some typ, .text n message, [], [], [], []⟩) =
      .ok (ascii "json: error calling MarshalJSON for type " ++ typ.text ++ [58, 32] ++ message) := by
  rfl

theorem marshalerErrorUnwrap_identity (e : CompatibilityError) :
    marshalerErrorUnwrap (some e) = .ok (nestedID e.nested) := rfl

theorem unmarshalTypeErrorUnwrap_identity (e : CompatibilityError) :
    unmarshalTypeErrorUnwrap (some e) = .ok (nestedID e.nested) := rfl

theorem unmarshalTypeError_missingType (e : CompatibilityError) (h : e.typ = none) :
    unmarshalTypeError (some e) = .nilPanic := by simp [unmarshalTypeError, h]

end GoLean.Encoding.Json
