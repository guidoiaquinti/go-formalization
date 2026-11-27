import GoLean.Strings.ReaderReset

namespace GoLean.Strings

/-- At cursor zero the receiver is unchanged. A positive cursor can be
    decremented even if it lies beyond the end of the backing string.
    Separate results represent an ordinary error and a nil receiver panic.
    https://github.com/golang/go/blob/go1.27.0/src/strings/reader.go#L76-L84 -/
def readerUnreadByte (receiver : Option ReaderState) :
    Option (Option String × ReaderState) := do
  let s ← receiver
  if s.cursor = 0 then
    return (some "strings.Reader.UnreadByte: at beginning of string", s)
  else
    let s := { s with previousRune := -1 }
    return (none, { s with cursor := s.cursor - 1 })

theorem readerUnreadByte_at_start (s : ReaderState) (h : s.cursor = 0) :
    readerUnreadByte (some s) =
      some (some "strings.Reader.UnreadByte: at beginning of string", s) := by
  change (if s.cursor = 0 then _ else _) = _
  rw [ite_eq_left h]
  rfl

theorem readerUnreadByte_positive (s : ReaderState) (h : 0 < s.cursor) :
    readerUnreadByte (some s) =
      some (none, ⟨s.contents, s.cursor - 1, -1⟩) := by
  change (if s.cursor = 0 then _ else _) = _
  rw [ite_eq_right (Nat.ne_of_gt h)]
  rfl

theorem readerUnreadByte_nil : readerUnreadByte none = none := rfl

end GoLean.Strings
