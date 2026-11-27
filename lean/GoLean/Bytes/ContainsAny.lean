import GoLean.Bytes.IndexAny

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:82–84. Report whether IndexAny
    returned a matching byte offset. -/
def containsAny (input chars : List (BitVec 8)) : Bool :=
  (indexAny input chars).isSome

theorem containsAny_spec (input chars : List (BitVec 8)) :
    containsAny input chars = true ↔ ∃ offset, indexAny input chars = some offset := by
  cases result : indexAny input chars <;> simp [containsAny, result]

theorem containsAny_emptyInput (chars : List (BitVec 8)) : containsAny [] chars = false := by
  simp [containsAny, indexAny_emptyInput]

theorem containsAny_emptyChars (input : List (BitVec 8)) : containsAny input [] = false := by
  simp [containsAny, indexAny_emptyChars]

end GoLean.Bytes
