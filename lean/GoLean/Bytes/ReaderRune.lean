import GoLean.Bytes.ReaderSetup
import GoLean.Unicode.Utf8.DecodeRune

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/reader.go:87–100. The exhausted path clears
    the marker. The live path saves the current cursor before decoding. -/
def readerReadRune (receiver : Option ReaderState) :
    Option (Nat × Nat × Bool × ReaderState) :=
  match receiver with
  | none => none
  | some state =>
    match ((state.bytes.getD []).drop state.cursor).map BitVec.toNat with
    | [] => some (0, 0, true, { state with previousRune := -1 })
    | head :: rest =>
      let (rune, width) := if head < 128 then (head, 1)
        else GoLean.Unicode.Utf8.decodeRune (head :: rest)
      some (rune, width, false,
        { state with cursor := state.cursor + width, previousRune := Int.ofNat state.cursor })

theorem readerReadRune_eof (s : ReaderState) (atEnd : (s.bytes.getD []).length ≤ s.cursor) :
    readerReadRune (some s) = some (0, 0, true, { s with previousRune := -1 }) := by
  dsimp only [readerReadRune]
  rw [List.drop_eq_nil_iff.mpr atEnd]
  rfl

theorem readerReadRune_live (s : ReaderState) (head : Nat) (rest : List Nat)
    (suffix : ((s.bytes.getD []).drop s.cursor).map BitVec.toNat = head :: rest) :
    readerReadRune (some s) =
      let (rune, width) := GoLean.Unicode.Utf8.decodeRune (head :: rest)
      some (rune, width, false,
        { s with cursor := s.cursor + width, previousRune := Int.ofNat s.cursor }) := by
  dsimp only [readerReadRune]
  rw [suffix]
  by_cases ascii : head < 128
  · simp [GoLean.Unicode.Utf8.decodeRune, ascii]
  · simp [ascii]

theorem readerReadRune_nil : readerReadRune none = none := rfl

/-- Go 1.27.0 src/bytes/reader.go:103–113. Check the cursor before the
    previous-rune marker. A successful rewind consumes the marker. -/
def readerUnreadRune (receiver : Option ReaderState) :
    Option (Option String × ReaderState) := do
  let state ← receiver
  if state.cursor = 0 then
    return (some "bytes.Reader.UnreadRune: at beginning of slice", state)
  else if state.previousRune < 0 then
    return (some "bytes.Reader.UnreadRune: previous operation was not ReadRune", state)
  else
    return (none, { state with cursor := state.previousRune.toNat, previousRune := -1 })

theorem readerUnreadRune_at_start (s : ReaderState) (h : s.cursor = 0) :
    readerUnreadRune (some s) =
      some (some "bytes.Reader.UnreadRune: at beginning of slice", s) := by
  change (if s.cursor = 0 then _ else _) = _
  rw [ite_eq_left h]
  rfl

theorem readerUnreadRune_missing_rune (s : ReaderState) (h : 0 < s.cursor)
    (missing : s.previousRune < 0) :
    readerUnreadRune (some s) =
      some (some "bytes.Reader.UnreadRune: previous operation was not ReadRune", s) := by
  change (if s.cursor = 0 then _ else if s.previousRune < 0 then _ else _) = _
  rw [ite_eq_right (Nat.ne_of_gt h), ite_eq_left missing]
  rfl

theorem readerUnreadRune_rewind (s : ReaderState) (h : 0 < s.cursor)
    (marked : 0 ≤ s.previousRune) :
    readerUnreadRune (some s) =
      some (none, { s with cursor := s.previousRune.toNat, previousRune := -1 }) := by
  change (if s.cursor = 0 then _ else if s.previousRune < 0 then _ else _) = _
  rw [ite_eq_right (Nat.ne_of_gt h), ite_eq_right (Int.not_lt.mpr marked)]
  rfl

theorem readerUnreadRune_nil : readerUnreadRune none = none := rfl

end GoLean.Bytes
