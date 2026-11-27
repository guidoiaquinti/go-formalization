import Init.Data.List.Lemmas

namespace GoLean.Bytes

/-- Go 1.27.0 Reader.Len for nonnegative seek positions, including
    positions past the underlying slice. None denotes a nil receiver panic. -/
def readerLen (receiver : Option (Nat × Nat)) : Option Nat :=
  receiver.map fun (length, offset) => length - offset

theorem readerLen_unread (bytes : List α) (offset : Nat) :
    readerLen (some (bytes.length, offset)) = some (bytes.drop offset).length := by
  simp [readerLen, List.length_drop]

theorem readerLen_pastEnd (length offset : Nat) (past : length ≤ offset) :
    readerLen (some (length, offset)) = some 0 := by
  simp [readerLen, Nat.sub_eq_zero_of_le past]

theorem readerLen_nil : readerLen none = none := rfl

end GoLean.Bytes
