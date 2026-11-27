import GoLean.Strings.ReaderReset

namespace GoLean.Strings

inductive ReadAtError where
  | negativeOffset
  | eof

/-- The offset check precedes the receiver access. An absent result
    represents a nil receiver panic. Reader state is read-only.
    https://github.com/golang/go/blob/go1.27.0/src/strings/reader.go#L49-L63 -/
def readerReadAt (receiver : Option ReaderState) (offset : Int) (dst : List (BitVec 8)) :
    Option (Nat × Option ReadAtError × List (BitVec 8)) :=
  if offset < 0 then some (0, some .negativeOffset, dst)
  else do
    let s ← receiver
    if offset.toNat ≥ s.contents.length then return (0, some .eof, dst)
    else
      let copied := (s.contents.drop offset.toNat).take dst.length
      let n := copied.length
      return (n, if n < dst.length then some .eof else none, copied ++ dst.drop n)

theorem readerReadAt_negative (r : Option ReaderState) (offset : Int) (h : offset < 0)
    (dst : List (BitVec 8)) :
    readerReadAt r offset dst = some (0, some .negativeOffset, dst) := by
  simp [readerReadAt, h]

theorem readerReadAt_nil (offset : Int) (h : 0 ≤ offset) (dst : List (BitVec 8)) :
    readerReadAt none offset dst = none := by
  simp [readerReadAt, Int.not_lt.mpr h]

theorem readerReadAt_eof (s : ReaderState) (offset : Int) (h : 0 ≤ offset)
    (atEnd : s.contents.length ≤ offset.toNat) (dst : List (BitVec 8)) :
    readerReadAt (some s) offset dst = some (0, some .eof, dst) := by
  unfold readerReadAt
  rw [ite_eq_right (Int.not_lt.mpr h)]
  change (if offset.toNat ≥ s.contents.length then _ else _) = _
  rw [ite_eq_left atEnd]
  rfl

theorem readerReadAt_copy (s : ReaderState) (offset : Int) (h : 0 ≤ offset)
    (inBounds : offset.toNat < s.contents.length) (dst : List (BitVec 8)) :
    readerReadAt (some s) offset dst =
      let n := min dst.length (s.contents.length - offset.toNat)
      some (n, if n < dst.length then some .eof else none,
        (s.contents.drop offset.toNat).take n ++ dst.drop n) := by
  unfold readerReadAt
  rw [ite_eq_right (Int.not_lt.mpr h)]
  change (if offset.toNat ≥ s.contents.length then _ else _) = _
  rw [ite_eq_right (Nat.not_le.mpr inBounds)]
  change some (_, _, _) = _
  simp only [List.length_take, List.length_drop]
  congr 3
  rw [← List.length_drop (l := s.contents) (i := offset.toNat)]
  rcases Nat.le_total dst.length (s.contents.drop offset.toNat).length with fits | short
  · rw [Nat.min_eq_left fits]
  · rw [Nat.min_eq_right short, List.take_of_length_le short, List.take_length]

end GoLean.Strings
