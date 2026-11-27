import GoLean.Bytes.ReaderSetup

namespace GoLean.Bytes

/-- A live receiver clears the rune marker, including at EOF. -/
def readerReadByte (receiver : Option ReaderState) :
    Option (BitVec 8 × Bool × ReaderState) := do
  let state ← receiver
  let state := { state with previousRune := -1 }
  let contents := state.bytes.getD []
  if h : state.cursor < contents.length then
    return (contents[state.cursor], false, { state with cursor := state.cursor + 1 })
  else
    return (0, true, state)

theorem readerReadByte_live (s : ReaderState) :
    readerReadByte (some s) =
      match (s.bytes.getD [])[s.cursor]? with
      | some byte => some (byte, false, { s with cursor := s.cursor + 1, previousRune := -1 })
      | none => some (0, true, { s with previousRune := -1 }) := by
  by_cases inBounds : s.cursor < (s.bytes.getD []).length
  · change (if h : s.cursor < (s.bytes.getD []).length then _ else _) = _
    rw [dite_eq_left inBounds, List.getElem?_eq_getElem inBounds]
    rfl
  · change (if h : s.cursor < (s.bytes.getD []).length then _ else _) = _
    rw [dite_eq_right inBounds, List.getElem?_eq_none (by omega)]
    rfl

theorem readerReadByte_nil : readerReadByte none = none := rfl

/-- At cursor zero the freshly constructed error leaves the receiver intact.
    Positive cursors may be decremented even beyond the end of the slice. -/
def readerUnreadByte (receiver : Option ReaderState) :
    Option (Option String × ReaderState) := do
  let state ← receiver
  if state.cursor = 0 then
    return (some "bytes.Reader.UnreadByte: at beginning of slice", state)
  else
    return (none, { state with cursor := state.cursor - 1, previousRune := -1 })

theorem readerUnreadByte_start (s : ReaderState) (h : s.cursor = 0) :
    readerUnreadByte (some s) =
      some (some "bytes.Reader.UnreadByte: at beginning of slice", s) := by
  change (if s.cursor = 0 then _ else _) = _
  rw [ite_eq_left h]
  rfl

theorem readerUnreadByte_positive (s : ReaderState) (h : 0 < s.cursor) :
    readerUnreadByte (some s) =
      some (none, { s with cursor := s.cursor - 1, previousRune := -1 }) := by
  change (if s.cursor = 0 then _ else _) = _
  rw [ite_eq_right (Nat.ne_of_gt h)]
  rfl

theorem readerUnreadByte_nil : readerUnreadByte none = none := rfl

end GoLean.Bytes
