import GoLean.Strings.BuilderReset

namespace GoLean.Strings

inductive BuilderIdentity where
  | unset
  | self
  | foreign

/-- Go 1.27.0 strings/builder.go:32–42. Addresses are represented by
    their relation to the receiver. The buffer is observed unchanged.
    None represents the nil receiver panic. -/
def builderCopyCheck (receiver : Option (BuilderIdentity × BuilderState)) :
    Option (Bool × BuilderIdentity × BuilderState) := do
  let (identity, b) ← receiver
  match identity with
  | .unset => return (false, .self, { b with identitySet := true })
  | .self => return (false, .self, b)
  | .foreign => return (true, .foreign, b)

theorem builderCopyCheck_unset (b : BuilderState) :
    builderCopyCheck (some (.unset, b)) =
      some (false, .self, { b with identitySet := true }) := rfl

theorem builderCopyCheck_self (b : BuilderState) :
    builderCopyCheck (some (.self, b)) = some (false, .self, b) := rfl

theorem builderCopyCheck_foreign (b : BuilderState) :
    builderCopyCheck (some (.foreign, b)) = some (true, .foreign, b) := rfl

theorem builderCopyCheck_nil :
    builderCopyCheck none = none := rfl

end GoLean.Strings
