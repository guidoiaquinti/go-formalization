import GoLean.Strings.IndexAny

namespace GoLean.Strings

/-- Go 1.27.0 strings/strings.go:67–69. IndexAny represents -1 as none. -/
def containsAny (input chars : List Nat) : Bool :=
  (indexAny input chars).isSome

theorem containsAny_spec (input chars : List Nat) :
    containsAny input chars = true ↔ ∃ offset, indexAny input chars = some offset := by
  cases h : indexAny input chars <;> simp [containsAny, h]

theorem containsAny_emptyChars (input : List Nat) : containsAny input [] = false := rfl

theorem containsAny_singleHighByte (input : List Nat) (byte : Nat) (high : 128 ≤ byte) :
    containsAny input [byte] = (indexRune input 0xfffd).isSome := by
  simp only [containsAny, indexAny_singleHighByte input byte high]

end GoLean.Strings
