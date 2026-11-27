import GoLean.Strings.BuilderReset

namespace GoLean.Strings

/-- A nil slice has zero capacity. A nil receiver panics.
    https://github.com/golang/go/blob/go1.27.0/src/strings/builder.go#L53-L56 -/
def builderCap (receiver : Option BuilderState) : Option Nat := do
  let state ← receiver
  let (_, cap) := state.buffer.getD ([], 0)
  return cap

theorem builderCap_live (state : BuilderState) :
    builderCap (some state) = some ((state.buffer.map Prod.snd).getD 0) := by
  cases state with
  | mk identity buffer => cases buffer <;> rfl

theorem builderCap_nil : builderCap none = none := rfl

end GoLean.Strings
