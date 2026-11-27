import GoLean.Strings.BuilderReset

namespace GoLean.Strings

inductive PrivateGrowStatus where
  | ok
  | allocationLength
  | sliceBounds

/-- Go 1.27.0 strings/builder.go:66–70. The allocation size is
    2*oldCap+n, and the rounded capacity controls reslicing. This helper
    does not run copyCheck. Allocation exhaustion is outside the model. -/
def builderGrowPrivate (receiver : Option BuilderState) (request : Int)
    (allocate : Nat → Nat) : Option (PrivateGrowStatus × BuilderState) := do
  let b ← receiver
  let (bytes, capacity) := b.buffer.getD ([], 0)
  let size := 2 * (capacity : Int) + request
  if size < 0 then return (.allocationLength, b)
  else
    let rounded := allocate size.toNat
    if rounded < bytes.length then return (.sliceBounds, b)
    else return (.ok, { b with buffer := some (bytes, rounded) })

theorem builderGrowPrivate_negativeSize (b : BuilderState) (n : Int) (allocate : Nat → Nat)
    (negative : 2 * ((b.buffer.getD ([], 0)).2 : Int) + n < 0) :
    builderGrowPrivate (some b) n allocate = some (.allocationLength, b) := by
  simp [builderGrowPrivate, negative]

theorem builderGrowPrivate_smallCapacity (b : BuilderState) (n : Int) (allocate : Nat → Nat)
    (valid : ¬ 2 * ((b.buffer.getD ([], 0)).2 : Int) + n < 0)
    (small : allocate (2 * ((b.buffer.getD ([], 0)).2 : Int) + n).toNat <
      (b.buffer.getD ([], 0)).1.length) :
    builderGrowPrivate (some b) n allocate = some (.sliceBounds, b) := by
  simp [builderGrowPrivate, valid, small]

theorem builderGrowPrivate_success (b : BuilderState) (n : Int) (allocate : Nat → Nat)
    (valid : ¬ 2 * ((b.buffer.getD ([], 0)).2 : Int) + n < 0)
    (fits : (b.buffer.getD ([], 0)).1.length ≤
      allocate (2 * ((b.buffer.getD ([], 0)).2 : Int) + n).toNat) :
    builderGrowPrivate (some b) n allocate =
      some (.ok, { b with buffer := some ((b.buffer.getD ([], 0)).1,
        allocate (2 * ((b.buffer.getD ([], 0)).2 : Int) + n).toNat) }) := by
  simp [builderGrowPrivate, valid, Nat.not_lt.mpr fits]

theorem builderGrowPrivate_nil (n : Int) (allocate : Nat → Nat) :
    builderGrowPrivate none n allocate = none := rfl

end GoLean.Strings
