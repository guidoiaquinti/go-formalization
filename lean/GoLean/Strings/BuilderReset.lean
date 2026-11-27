import Init.Data.BitVec.Basic

namespace GoLean.Strings

/-- Observable buffer contents and capacity, plus whether a copy-check identity is set. -/
structure BuilderState where
  identitySet : Bool
  buffer : Option (List (BitVec 8) × Nat)

/-- A nil receiver panics on the first field assignment. A live receiver has
    its copy-check identity and slice cleared, in that order.
    https://github.com/golang/go/blob/go1.27.0/src/strings/builder.go#L59-L63 -/
def builderReset (receiver : Option BuilderState) : Option BuilderState := do
  let state ← receiver
  let state := { state with identitySet := false }
  return { state with buffer := none }

theorem builderReset_live (state : BuilderState) :
    builderReset (some state) = some ⟨false, none⟩ := rfl

theorem builderReset_nil : builderReset none = none := rfl

end GoLean.Strings
