namespace GoLean.Internal.Dag

/-- Go 1.27.0 converts the defined string type to its underlying string type.
    The model uses byte lists because Go strings may contain arbitrary bytes.
    https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L188-L192 -/
structure SyntaxErrorValue where
  bytes : List UInt8

def syntaxErrorError (receiver : SyntaxErrorValue) : List UInt8 :=
  receiver.bytes

theorem syntaxErrorError_spec (input : List UInt8) :
    syntaxErrorError ⟨input⟩ = input := rfl

end GoLean.Internal.Dag
