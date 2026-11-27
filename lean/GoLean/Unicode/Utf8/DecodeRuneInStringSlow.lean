import GoLean.Unicode.Utf8.DecodeRuneInString

namespace GoLean.Unicode.Utf8

/-- Go 1.27.0 src/unicode/utf8/utf8.go:231–269. The shared model captures
    the table's ASCII and illegal-lead mask as the corresponding branch. -/
def decodeRuneInStringSlow (bytes : List Nat) : Nat × Nat :=
  match bytes with
  | [] => (0xfffd, 0)
  | lead :: rest => decodeNonempty lead rest

theorem decodeRuneInStringSlow_spec (bytes : List Nat) :
    decodeRuneInStringSlow bytes = decodeRuneInString bytes := by
  cases bytes <;> rfl

theorem decodeRuneInStringSlow_ascii (lead : Nat) (rest : List Nat) (ascii : lead < 128) :
    decodeRuneInStringSlow (lead :: rest) = (lead, 1) := by
  rw [decodeRuneInStringSlow_spec]
  exact decodeRuneInString_ascii lead rest ascii

end GoLean.Unicode.Utf8
