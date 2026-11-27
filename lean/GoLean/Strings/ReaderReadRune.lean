import GoLean.Unicode.Utf8.DecodeRuneInString
import GoLean.Strings.ReaderReset

namespace GoLean.Strings

abbrev readerRune := GoLean.Unicode.Utf8.decodeNonempty

/-- The empty suffix represents EOF. A nonempty suffix records the
    original cursor and advances by the decoded width.
    https://github.com/golang/go/blob/go1.27.0/src/strings/reader.go#L86-L96 -/
def readerReadRune (receiver : Option ReaderState) :
    Option (Nat × Nat × Bool × ReaderState) := do
  let s ← receiver
  match (s.contents.drop s.cursor).map BitVec.toNat with
  | [] => return (0, 0, true, { s with previousRune := -1 })
  | lead :: rest =>
    let s := { s with previousRune := s.cursor }
    let (rune, width) := readerRune lead rest
    return (rune, width, false, { s with cursor := s.cursor + width })

theorem readerReadRune_eof (s : ReaderState) (atEnd : s.contents.length ≤ s.cursor) :
    readerReadRune (some s) = some (0, 0, true, { s with previousRune := -1 }) := by
  have empty : s.contents.drop s.cursor = [] := List.drop_eq_nil_iff.mpr atEnd
  change (match (s.contents.drop s.cursor).map BitVec.toNat with
    | [] => _
    | _ :: _ => _) = _
  rw [empty]
  rfl

theorem readerReadRune_live (s : ReaderState) (lead : Nat) (rest : List Nat)
    (suffix : (s.contents.drop s.cursor).map BitVec.toNat = lead :: rest) :
    readerReadRune (some s) =
      let (rune, width) := readerRune lead rest
      some (rune, width, false, ⟨s.contents, s.cursor + width, s.cursor⟩) := by
  change (match (s.contents.drop s.cursor).map BitVec.toNat with
    | [] => _
    | _ :: _ => _) = _
  rw [suffix]
  rfl

theorem readerReadRune_nil : readerReadRune none = none := rfl

end GoLean.Strings
