import GoLean.Bytes.ReaderSetup
import Init.Data.BitVec.Lemmas

namespace GoLean.Bytes

inductive SeekError where
  | invalidWhence
  | negativePosition

private def seekCandidate (s : ReaderState) (offset : BitVec 64) (whence : Int) :
    Option (BitVec 64) :=
  match whence with
  | 0 => some offset
  | 1 => some (BitVec.ofNat 64 s.cursor + offset)
  | 2 => some (BitVec.ofNat 64 (s.bytes.getD []).length + offset)
  | _ => none

/-- The prior rune marker is cleared before all whence and position checks.
    Adding int64 offsets wraps before the resulting position is tested.
    https://github.com/golang/go/blob/go1.27.0/src/bytes/reader.go#L116-L134 -/
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
      next.bytes = s.bytes ∧ next.spare = s.spare ∧ next.previousRune = -1 ∧
        match err with
        | none => 0 ≤ pos ∧ next.cursor = pos.toNat
        | some _ => pos = 0 ∧ next.cursor = s.cursor := by
  change (match (match seekCandidate s offset whence with
    | none => some _
    | some position => if position.toInt < 0 then some _ else some _ : Option (Int × Option SeekError × ReaderState)) with
    | none => False
    | some (pos, err, next) =>
      next.bytes = s.bytes ∧ next.spare = s.spare ∧ next.previousRune = -1 ∧
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

inductive ReaderWriterResponse where
  | returned (count : Int) (failed : Bool)
  | panicked
  | missing

inductive ReaderWriteStatus where
  | ok
  | writerError
  | shortWrite
  | writerPanic
  | invalidCount
  | nilWriter

structure ReaderWriteCall where
  bytes : List (BitVec 8)
  capacity : Nat

structure ReaderWriteResult where
  count : Int
  status : ReaderWriteStatus
  cursor : Int
  previousRune : Int
  call : Option ReaderWriteCall

private def writerReceipt (s : ReaderState) : ReaderWriteCall :=
  let suffix := (s.bytes.getD []).drop s.cursor
  ⟨suffix, suffix.length + s.spare.length⟩

/-- Go 1.27.0 src/bytes/reader.go:137–153. The Write call sees the
    original unread slice after clearing the rune marker. A malformed
    count may move the cursor below zero. -/
def readerWriteTo (receiver : Option ReaderState) (response : ReaderWriterResponse) :
    Option ReaderWriteResult := do
  let s ← receiver
  if s.cursor ≥ (s.bytes.getD []).length then
    return ⟨0, .ok, s.cursor, -1, none⟩
  else
    let receipt := writerReceipt s
    match response with
    | .missing => return ⟨0, .nilWriter, s.cursor, -1, none⟩
    | .panicked => return ⟨0, .writerPanic, s.cursor, -1, some receipt⟩
    | .returned n failed =>
      if n > (receipt.bytes.length : Int) then
        return ⟨0, .invalidCount, s.cursor, -1, some receipt⟩
      else
        let status := if failed then .writerError else if n ≠ receipt.bytes.length then .shortWrite else .ok
        return ⟨n, status, (s.cursor : Int) + n, -1, some receipt⟩

theorem readerWriteTo_eof (s : ReaderState) (atEnd : (s.bytes.getD []).length ≤ s.cursor)
    (response : ReaderWriterResponse) :
    readerWriteTo (some s) response = some ⟨0, .ok, s.cursor, -1, none⟩ := by
  change (if s.cursor ≥ (s.bytes.getD []).length then _ else _) = _
  rw [ite_eq_left atEnd]
  rfl

theorem readerWriteTo_count (s : ReaderState) (inBounds : s.cursor < (s.bytes.getD []).length)
    (n : Int) (failed : Bool) :
    readerWriteTo (some s) (.returned n failed) =
      let suffix := (s.bytes.getD []).drop s.cursor
      let call := some (ReaderWriteCall.mk suffix (suffix.length + s.spare.length))
      if n > (suffix.length : Int) then
        some ⟨0, .invalidCount, s.cursor, -1, call⟩
      else
        let status := if failed then .writerError else if n ≠ suffix.length then .shortWrite else .ok
        some ⟨n, status, (s.cursor : Int) + n, -1, call⟩ := by
  change (if s.cursor ≥ (s.bytes.getD []).length then _ else _) = _
  rw [ite_eq_right (Nat.not_le.mpr inBounds)]
  rfl

theorem readerWriteTo_panic (s : ReaderState) (inBounds : s.cursor < (s.bytes.getD []).length) :
    readerWriteTo (some s) .panicked =
      some ⟨0, .writerPanic, s.cursor, -1, some (writerReceipt s)⟩ := by
  change (if s.cursor ≥ (s.bytes.getD []).length then _ else _) = _
  rw [ite_eq_right (Nat.not_le.mpr inBounds)]
  rfl

theorem readerWriteTo_nil (response : ReaderWriterResponse) :
    readerWriteTo none response = none := rfl

end GoLean.Bytes
