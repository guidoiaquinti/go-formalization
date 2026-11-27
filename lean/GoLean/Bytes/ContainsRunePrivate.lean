import GoLean.Unicode.Utf8.DecodeRuneInString

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:976–983. Traverse the cutset as raw
    string bytes using Go's range decoding and stop at the first matching rune. -/
def containsRunePrivate : List Nat → Int → Bool
  | [], _ => false
  | lead :: rest, seek =>
    let (decoded, width) := GoLean.Unicode.Utf8.decodeRuneInString (lead :: rest)
    if Int.ofNat decoded == seek then true
    else containsRunePrivate (rest.drop (width - 1)) seek
termination_by input _ => input.length
decreasing_by simp [List.length_drop]; omega

theorem containsRunePrivate_empty (seek : Int) :
    containsRunePrivate [] seek = false := by
  rw [containsRunePrivate.eq_1]

theorem containsRunePrivate_cons (lead : Nat) (rest : List Nat) (seek : Int) :
    let (decoded, width) := GoLean.Unicode.Utf8.decodeRuneInString (lead :: rest)
    containsRunePrivate (lead :: rest) seek =
      if Int.ofNat decoded == seek then true else containsRunePrivate (rest.drop (width - 1)) seek := by
  rw [containsRunePrivate.eq_2]

theorem containsRunePrivate_negative (input : List Nat) (seek : Int) (negative : seek < 0) :
    containsRunePrivate input seek = false := by
  cases input with
  | nil => exact containsRunePrivate_empty seek
  | cons lead rest =>
    rw [containsRunePrivate.eq_2]
    simp only
    split
    · rename_i equal
      simp at equal
      omega
    · exact containsRunePrivate_negative _ seek negative
termination_by input.length
decreasing_by simp [List.length_drop]; omega

end GoLean.Bytes
