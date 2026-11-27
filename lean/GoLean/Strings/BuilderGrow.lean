import GoLean.Strings.BuilderReset

namespace GoLean.Strings

inductive BuilderGrowStatus where
  | ok
  | copied
  | negative

/-- Go 1.27.0 strings/builder.go:66–82. An allocator parameter represents
    MakeNoZero rounding for a positive, representable allocation request. -/
def builderGrow (receiver : Option BuilderState) (copied : Bool) (request : Int)
    (allocate : Nat → Nat) : Option (BuilderGrowStatus × BuilderState) := do
  let b ← receiver
  if b.identitySet && copied then
    return (.copied, b)
  else
    let b := { b with identitySet := true }
    if request < 0 then
      return (.negative, b)
    else
      let (bytes, capacity) := b.buffer.getD ([], 0)
      if capacity - bytes.length < request.toNat then
        return (.ok, { b with buffer := some (bytes, allocate (2 * capacity + request.toNat)) })
      else
        return (.ok, b)

theorem builderGrow_copied (b : BuilderState) (request : Int) (allocate : Nat → Nat)
    (identified : b.identitySet = true) :
    builderGrow (some b) true request allocate = some (.copied, b) := by
  simp [builderGrow, identified]

theorem builderGrow_negative (b : BuilderState) (request : Int) (allocate : Nat → Nat)
    (negative : request < 0) :
    builderGrow (some b) false request allocate =
      some (.negative, { b with identitySet := true }) := by
  simp [builderGrow, negative]

theorem builderGrow_space (b : BuilderState) (request : Nat) (allocate : Nat → Nat)
    (fits : request ≤ (b.buffer.getD ([], 0)).2 - (b.buffer.getD ([], 0)).1.length) :
    builderGrow (some b) false request allocate = some (.ok, { b with identitySet := true }) := by
  simp [builderGrow, Nat.not_lt.mpr fits]

theorem builderGrow_allocate (b : BuilderState) (request : Nat) (allocate : Nat → Nat)
    (grows : (b.buffer.getD ([], 0)).2 - (b.buffer.getD ([], 0)).1.length < request) :
    builderGrow (some b) false request allocate =
      some (.ok, ⟨true, some ((b.buffer.getD ([], 0)).1,
        allocate (2 * (b.buffer.getD ([], 0)).2 + request))⟩) := by
  simp [builderGrow, grows]

end GoLean.Strings
