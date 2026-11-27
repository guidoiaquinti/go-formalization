import GoLean.Strings.ReaderReset

namespace GoLean.Strings

/-- Check the cursor before the rune marker. Only a successful rewind
    modifies the state. An absent result represents a nil receiver panic.
    https://github.com/golang/go/blob/go1.27.0/src/strings/reader.go#L98-L109 -/
def readerUnreadRune (receiver : Option ReaderState) :
    Option (Option String × ReaderState) := do
  let s ← receiver
  if s.cursor = 0 then
    return (some "strings.Reader.UnreadRune: at beginning of string", s)
  else if s.previousRune < 0 then
    return (some "strings.Reader.UnreadRune: previous operation was not ReadRune", s)
  else
    let s := { s with cursor := s.previousRune.toNat }
    return (none, { s with previousRune := -1 })

theorem readerUnreadRune_at_start (s : ReaderState) (h : s.cursor = 0) :
    readerUnreadRune (some s) =
      some (some "strings.Reader.UnreadRune: at beginning of string", s) := by
  change (if s.cursor = 0 then _ else _) = _
  rw [ite_eq_left h]
  rfl

theorem readerUnreadRune_missing_rune (s : ReaderState) (h : 0 < s.cursor)
    (missing : s.previousRune < 0) :
    readerUnreadRune (some s) =
      some (some "strings.Reader.UnreadRune: previous operation was not ReadRune", s) := by
  change (if s.cursor = 0 then _ else if s.previousRune < 0 then _ else _) = _
  rw [ite_eq_right (Nat.ne_of_gt h), ite_eq_left missing]
  rfl

theorem readerUnreadRune_rewind (s : ReaderState) (h : 0 < s.cursor)
    (marked : 0 ≤ s.previousRune) :
    readerUnreadRune (some s) =
      some (none, ⟨s.contents, s.previousRune.toNat, -1⟩) := by
  change (if s.cursor = 0 then _ else if s.previousRune < 0 then _ else _) = _
  rw [ite_eq_right (Nat.ne_of_gt h), ite_eq_right (Int.not_lt.mpr marked)]
  rfl

theorem readerUnreadRune_nil : readerUnreadRune none = none := rfl

end GoLean.Strings
