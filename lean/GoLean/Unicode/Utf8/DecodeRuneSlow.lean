import GoLean.Unicode.Utf8.DecodeRuneInStringSlow

namespace GoLean.Unicode.Utf8

/-- Go 1.27.0 src/unicode/utf8/utf8.go:171–208. The slice helper has
    the same acceptance table, ASCII mask and payload operations as the
    string helper. Nil and empty slice contents are represented by []. -/
def decodeRuneSlow (bytes : List Nat) : Nat × Nat :=
  match bytes with
  | [] => (0xfffd, 0)
  | lead :: rest => decodeNonempty lead rest

theorem decodeRuneSlow_spec (bytes : List Nat) :
    decodeRuneSlow bytes = decodeRuneInStringSlow bytes := by
  cases bytes <;> rfl

theorem decodeRuneSlow_ascii (lead : Nat) (rest : List Nat) (ascii : lead < 128) :
    decodeRuneSlow (lead :: rest) = (lead, 1) := by
  rw [decodeRuneSlow_spec]
  exact decodeRuneInStringSlow_ascii lead rest ascii

end GoLean.Unicode.Utf8
