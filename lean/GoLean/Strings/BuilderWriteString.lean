import GoLean.Strings.BuilderReset

namespace GoLean.Strings

structure BuilderWriteResult where
  count : Option Nat
  contents : List (BitVec 8)
  identitySet : Bool

/-- The copy check runs before append. A missing receiver denotes a nil
    panic, and a missing count records a copy-check panic on a live builder.
    The result abstracts allocation capacity after append.
    https://github.com/golang/go/blob/go1.27.0/src/strings/builder.go#L32-L42
    https://github.com/golang/go/blob/go1.27.0/src/strings/builder.go#L110-L116 -/
def builderWriteString (receiver : Option BuilderState) (copied : Bool)
    (input : List (BitVec 8)) : Option BuilderWriteResult := do
  let b ← receiver
  let (contents, _) := b.buffer.getD ([], 0)
  if b.identitySet && copied then
    return ⟨none, contents, b.identitySet⟩
  else
    let b := { b with identitySet := true }
    return ⟨some input.length, contents ++ input, b.identitySet⟩

theorem builderWriteString_success (b : BuilderState) (copied : Bool)
    (input : List (BitVec 8)) (canWrite : (b.identitySet && copied) = false) :
    builderWriteString (some b) copied input =
      some ⟨some input.length, (b.buffer.getD ([], 0)).1 ++ input, true⟩ := by
  change (if b.identitySet && copied then _ else _) = _
  simp [canWrite, Option.getD]

theorem builderWriteString_copied (b : BuilderState) (input : List (BitVec 8))
    (identified : b.identitySet = true) :
    builderWriteString (some b) true input =
      some ⟨none, (b.buffer.getD ([], 0)).1, true⟩ := by
  change (if b.identitySet && true then _ else _) = _
  simp [identified, Option.getD]

theorem builderWriteString_nil (copied : Bool) (input : List (BitVec 8)) :
    builderWriteString none copied input = none := rfl

end GoLean.Strings
