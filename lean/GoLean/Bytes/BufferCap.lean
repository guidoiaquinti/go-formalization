import Init.Omega

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/buffer.go:100–101. The observable capacity
    belongs to the underlying slice. None is a nil receiver panic. -/
def bufferCap (receiver : Option Nat) : Option Nat :=
  receiver

theorem bufferCap_backing (capacity : Nat) :
    bufferCap (some capacity) = some capacity := rfl

theorem bufferCap_unreadGeometry (capacity offset : Nat) (reachable : offset ≤ capacity) :
    bufferCap (some capacity) = some (capacity - offset + offset) := by
  simp only [bufferCap, Option.some.injEq]
  omega

theorem bufferCap_nil : bufferCap none = none := rfl

end GoLean.Bytes
