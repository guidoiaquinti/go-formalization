import GoLean.Unicode.Utf8.DecodeRuneInString

namespace GoLean.Unicode.Utf8

/-- The slice decoder has a direct first-byte ASCII path and otherwise uses the
    same acceptance table and payload assembly as the string decoder.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf8/utf8.go#L157-L208 -/
def decodeRune : List Nat → Nat × Nat
  | [] => (0xfffd, 0)
  | lead :: rest => if lead < 128 then (lead, 1) else decodeNonempty lead rest

theorem decodeRune_spec (bytes : List Nat) :
    decodeRune bytes = decodeRuneInString bytes := by
  cases bytes with
  | nil => rfl
  | cons lead rest =>
    by_cases ascii : lead < 128
    · simp [decodeRune, ascii, decodeRuneInString_ascii lead rest ascii]
    · simp [decodeRune, decodeRuneInString, ascii]

end GoLean.Unicode.Utf8
