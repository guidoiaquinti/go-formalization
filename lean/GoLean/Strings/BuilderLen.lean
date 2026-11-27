import GoLean.Strings.BuilderReset

namespace GoLean.Strings

/-- The builder's byte length is independent of the allocation capacity.
    A missing receiver models a nil receiver panic.
    https://github.com/golang/go/blob/go1.27.0/src/strings/builder.go#L50-L51 -/
def builderLen (receiver : Option BuilderState) : Option Nat := do
  let state ← receiver
  let (bytes, _) := state.buffer.getD ([], 0)
  return bytes.length

theorem builderLen_live (state : BuilderState) :
    builderLen (some state) = some ((state.buffer.map Prod.fst).getD []).length := by
  cases state with
  | mk identity buffer => cases buffer <;> rfl

theorem builderLen_nil : builderLen none = none := rfl

end GoLean.Strings
