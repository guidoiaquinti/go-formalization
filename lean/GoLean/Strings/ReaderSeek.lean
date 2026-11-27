import GoLean.Strings.ReaderReset
import Init.Data.BitVec.Lemmas

namespace GoLean.Strings

inductive SeekError where
  | invalidWhence
  | negativePosition

private def seekCandidate (s : ReaderState) (offset : BitVec 64) (whence : Int) :
    Option (BitVec 64) :=
  match whence with
  | 0 => some offset
  | 1 => some (BitVec.ofNat 64 s.cursor + offset)
  | 2 => some (BitVec.ofNat 64 s.contents.length + offset)
  | _ => none

/-- The prior rune marker is cleared before all whence and position checks.
    Adding int64 offsets wraps before the resulting position is tested.
    https://github.com/golang/go/blob/go1.27.0/src/strings/reader.go#L111-L131 -/
def readerSeek (r : Option ReaderState) (offset : BitVec 64) (whence : Int) :
    Option (Int × Option SeekError × ReaderState) := do
  let s ← r
  let s := { s with previousRune := -1 }
  match seekCandidate s offset whence with
  | none => return (0, some .invalidWhence, s)
  | some position =>
    if position.toInt < 0 then
      return (0, some .negativePosition, s)
    else
      return (position.toInt, none, { s with cursor := position.toNat })

theorem readerSeek_live (s : ReaderState) (offset : BitVec 64) (whence : Int) :
    match readerSeek (some s) offset whence with
    | none => False
    | some (pos, err, next) =>
      next.contents = s.contents ∧ next.previousRune = -1 ∧
        match err with
        | none => 0 ≤ pos ∧ next.cursor = pos.toNat
        | some _ => pos = 0 ∧ next.cursor = s.cursor := by
  change (match (match seekCandidate s offset whence with
    | none => some _
    | some position => if position.toInt < 0 then some _ else some _ : Option (Int × Option SeekError × ReaderState)) with
    | none => False
    | some (pos, err, next) =>
      next.contents = s.contents ∧ next.previousRune = -1 ∧
        match err with
        | none => 0 ≤ pos ∧ next.cursor = pos.toNat
        | some _ => pos = 0 ∧ next.cursor = s.cursor)
  cases result : seekCandidate s offset whence with
  | none => simp
  | some position =>
    by_cases negative : position.toInt < 0
    · simp [negative]
    · have nonnegative : 0 ≤ position.toInt := Int.not_lt.mp negative
      have signed : BitVec.sle 0#64 position := by simpa [BitVec.sle] using nonnegative
      simp [negative, nonnegative, BitVec.toNat_toInt_of_sle signed]

theorem readerSeek_nil (offset : BitVec 64) (whence : Int) :
    readerSeek none offset whence = none := rfl

end GoLean.Strings
