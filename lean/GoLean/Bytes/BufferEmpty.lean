namespace GoLean.Bytes

/-- Go 1.27.0 Buffer.empty compares the stored length with its signed
    read cursor. None represents a nil receiver panic. -/
def bufferEmpty (receiver : Option (Nat × Int)) : Option Bool :=
  receiver.map fun (length, offset) => decide (Int.ofNat length ≤ offset)

theorem bufferEmpty_nil : bufferEmpty none = none := rfl

theorem bufferEmpty_end (length : Nat) :
    bufferEmpty (some (length, Int.ofNat length)) = some true := by
  simp [bufferEmpty]

theorem bufferEmpty_hasUnread (length : Nat) (offset : Int)
    (before : offset < Int.ofNat length) :
    bufferEmpty (some (length, offset)) = some false := by
  have outside : ¬Int.ofNat length ≤ offset := by omega
  simp only [bufferEmpty, Option.map_some, decide_eq_false outside]

end GoLean.Bytes
