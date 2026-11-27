import GoLean.Strings.ReaderReset

namespace GoLean.Strings

inductive WriterResult where
  | returned (n : Int) (failed : Bool)
  | panicked
  | missing

inductive WriteToStatus where
  | ok
  | writerError
  | shortWrite
  | writerPanic
  | invalidCount
  | nilWriter

structure WriteToResult where
  count : Int
  status : WriteToStatus
  cursor : Int
  previousRune : Int
  call : Option (Bool × List (BitVec 8))

/-- Record the selected writer method and its argument. A malformed writer
    can return a negative count, so the output cursor is signed.
    https://github.com/golang/go/blob/go1.27.0/src/strings/reader.go#L132-L149
    https://github.com/golang/go/blob/go1.27.0/src/io/io.go -/
def readerWriteTo (receiver : Option ReaderState) (stringWriter : Bool)
    (result : WriterResult) : Option WriteToResult := do
  let s ← receiver
  let s := { s with previousRune := -1 }
  if s.cursor ≥ s.contents.length then
    return ⟨0, .ok, s.cursor, s.previousRune, none⟩
  else
    let suffix := s.contents.drop s.cursor
    let call := some (stringWriter, suffix)
    match result with
    | .missing => return ⟨0, .nilWriter, s.cursor, s.previousRune, none⟩
    | .panicked => return ⟨0, .writerPanic, s.cursor, s.previousRune, call⟩
    | .returned n failed =>
      if n > (suffix.length : Int) then
        return ⟨0, .invalidCount, s.cursor, s.previousRune, call⟩
      else
        let status := if failed then .writerError else if n ≠ suffix.length then .shortWrite else .ok
        return ⟨n, status, (s.cursor : Int) + n, s.previousRune, call⟩

theorem readerWriteTo_eof (s : ReaderState) (atEnd : s.contents.length ≤ s.cursor)
    (method : Bool) (response : WriterResult) :
    readerWriteTo (some s) method response = some ⟨0, .ok, s.cursor, -1, none⟩ := by
  change (if s.cursor ≥ s.contents.length then _ else _) = _
  rw [ite_eq_left atEnd]
  rfl

theorem readerWriteTo_count (s : ReaderState) (inBounds : s.cursor < s.contents.length)
    (method : Bool) (n : Int) (failed : Bool) :
    readerWriteTo (some s) method (.returned n failed) =
      let suffix := s.contents.drop s.cursor
      if n > (suffix.length : Int) then
        some ⟨0, .invalidCount, s.cursor, -1, some (method, suffix)⟩
      else
        let status := if failed then .writerError else if n ≠ suffix.length then .shortWrite else .ok
        some ⟨n, status, (s.cursor : Int) + n, -1, some (method, suffix)⟩ := by
  change (if s.cursor ≥ s.contents.length then _ else _) = _
  rw [ite_eq_right (Nat.not_le.mpr inBounds)]
  rfl

theorem readerWriteTo_panic (s : ReaderState) (inBounds : s.cursor < s.contents.length)
    (method : Bool) :
    readerWriteTo (some s) method .panicked =
      some ⟨0, .writerPanic, s.cursor, -1, some (method, s.contents.drop s.cursor)⟩ := by
  change (if s.cursor ≥ s.contents.length then _ else _) = _
  rw [ite_eq_right (Nat.not_le.mpr inBounds)]
  rfl

theorem readerWriteTo_missing (s : ReaderState) (inBounds : s.cursor < s.contents.length)
    (method : Bool) :
    readerWriteTo (some s) method .missing = some ⟨0, .nilWriter, s.cursor, -1, none⟩ := by
  change (if s.cursor ≥ s.contents.length then _ else _) = _
  rw [ite_eq_right (Nat.not_le.mpr inBounds)]
  rfl

theorem readerWriteTo_nil (method : Bool) (response : WriterResult) :
    readerWriteTo none method response = none := rfl

end GoLean.Strings
