import GoLean.Strings.ReaderReset

namespace GoLean.Strings

/-- A live receiver clears the rune marker even when the cursor is at EOF.
    The Boolean result marks EOF, and an absent state represents a nil
    receiver panic. Public cursor positions are nonnegative.
    https://github.com/golang/go/blob/go1.27.0/src/strings/reader.go#L65-L74 -/
def readerReadByte (receiver : Option ReaderState) :
    Option (BitVec 8 × Bool × ReaderState) := do
  let s ← receiver
  let s := { s with previousRune := -1 }
  if h : s.cursor < s.contents.length then
    return (s.contents[s.cursor], false, { s with cursor := s.cursor + 1 })
  else
    return (0, true, s)

theorem readerReadByte_live (s : ReaderState) :
    readerReadByte (some s) =
      match s.contents[s.cursor]? with
      | some byte => some (byte, false, { s with cursor := s.cursor + 1, previousRune := -1 })
      | none => some (0, true, { s with previousRune := -1 }) := by
  by_cases inBounds : s.cursor < s.contents.length
  · change (if h : s.cursor < s.contents.length then _ else _) = _
    rw [dite_eq_left inBounds, List.getElem?_eq_getElem inBounds]
    rfl
  · change (if h : s.cursor < s.contents.length then _ else _) = _
    rw [dite_eq_right inBounds, List.getElem?_eq_none (by omega)]
    rfl

theorem readerReadByte_nil : readerReadByte none = none := rfl

end GoLean.Strings
