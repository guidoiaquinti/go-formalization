import GoLean.Strings.ReaderReset

namespace GoLean.Strings

/-- EOF is checked before clearing the rune marker. A successful short
    copy still returns a nil error. None represents a nil receiver panic.
    https://github.com/golang/go/blob/go1.27.0/src/strings/reader.go#L38-L47 -/
def readerRead (receiver : Option ReaderState) (dst : List (BitVec 8)) :
    Option (Nat × Bool × List (BitVec 8) × ReaderState) := do
  let s ← receiver
  if s.cursor ≥ s.contents.length then
    return (0, true, dst, s)
  else
    let s := { s with previousRune := -1 }
    let copied := (s.contents.drop s.cursor).take dst.length
    let n := copied.length
    return (n, false, copied ++ dst.drop n, { s with cursor := s.cursor + n })

theorem readerRead_eof (s : ReaderState) (atEnd : s.contents.length ≤ s.cursor)
    (dst : List (BitVec 8)) :
    readerRead (some s) dst = some (0, true, dst, s) := by
  change (if s.cursor ≥ s.contents.length then _ else _) = _
  rw [ite_eq_left atEnd]
  rfl

theorem readerRead_copy (s : ReaderState) (inBounds : s.cursor < s.contents.length)
    (dst : List (BitVec 8)) :
    readerRead (some s) dst =
      let n := min dst.length (s.contents.length - s.cursor)
      some (n, false, (s.contents.drop s.cursor).take dst.length ++ dst.drop n,
        ⟨s.contents, s.cursor + n, -1⟩) := by
  change (if s.cursor ≥ s.contents.length then _ else _) = _
  rw [ite_eq_right (Nat.not_le.mpr inBounds)]
  change some (_, _, _, _) = _
  simp only [List.length_take, List.length_drop]

theorem readerRead_nil (dst : List (BitVec 8)) : readerRead none dst = none := rfl

end GoLean.Strings
