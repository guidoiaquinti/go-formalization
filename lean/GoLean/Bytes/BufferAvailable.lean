import Init.Omega

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/buffer.go:103–104. Available is the space to
    append at the end of the underlying slice. None denotes a nil
    receiver panic. Reachable geometry has length no greater than capacity. -/
def bufferAvailable (receiver : Option (Nat × Nat)) : Option Nat :=
  receiver.map fun (length, capacity) => capacity - length

theorem bufferAvailable_unreadGeometry (length capacity offset : Nat)
    (reachable : offset ≤ length) (fits : length ≤ capacity) :
    bufferAvailable (some (length, capacity)) =
      some ((capacity - offset) - (length - offset)) := by
  simp only [bufferAvailable, Option.map_some, Option.some.injEq]
  omega

theorem bufferAvailable_full (length : Nat) :
    bufferAvailable (some (length, length)) = some 0 := by
  simp [bufferAvailable]

theorem bufferAvailable_nil : bufferAvailable none = none := rfl

end GoLean.Bytes
