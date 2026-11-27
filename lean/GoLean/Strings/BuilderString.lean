import GoLean.Strings.BuilderReset

namespace GoLean.Strings

/-- The builder exposes only the accumulated byte range, regardless of reserve.
    A missing receiver models a nil receiver panic.
    https://github.com/golang/go/blob/go1.27.0/src/strings/builder.go#L45-L48 -/
def builderString (receiver : Option BuilderState) : Option (List (BitVec 8)) := do
  let state ← receiver
  let (bytes, _) := state.buffer.getD ([], 0)
  return bytes

theorem builderString_live (state : BuilderState) :
    builderString (some state) = some ((state.buffer.map Prod.fst).getD []) := by
  cases state with
  | mk identity buffer => cases buffer <;> rfl

theorem builderString_nil : builderString none = none := rfl

end GoLean.Strings
