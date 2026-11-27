import GoLean.Bytes.ReaderSetup

namespace GoLean.Bytes

inductive ReaderDestination where
  | separate (contents : List (BitVec 8))
  | shared (offset length : Nat)

inductive ReaderCopyStatus where
  | ok
  | eof
  | negativeOffset

structure ReaderCopyResult where
  count : Nat
  status : ReaderCopyStatus
  destination : List (BitVec 8)
  state : Option ReaderState

private def readerStorage (state : Option ReaderState) : List (BitVec 8) :=
  state.map (fun s => s.bytes.getD [] ++ s.spare) |>.getD []

private def destinationContents (state : Option ReaderState) : ReaderDestination → List (BitVec 8)
  | .separate contents => contents
  | .shared offset length => ((readerStorage state).drop offset).take length

private def copyReaderBytes (state : ReaderState) (source : Nat) (destination : ReaderDestination) :
    Nat × List (BitVec 8) × ReaderState :=
  let old := destinationContents (some state) destination
  let copied := ((state.bytes.getD []).drop source).take old.length
  let visible := copied ++ old.drop copied.length
  let updated := match destination with
    | .separate _ => state
    | .shared offset _ =>
      let storage := readerStorage (some state)
      let patched := storage.take offset ++ copied ++ storage.drop (offset + copied.length)
      let size := (state.bytes.getD []).length
      { state with bytes := state.bytes.map (fun _ => patched.take size), spare := patched.drop size }
  (copied.length, visible, updated)

/-- Go 1.27.0 src/bytes/reader.go:39–47. EOF precedes clearing the marker.
    Copy snapshots the source before writing to possibly overlapping storage. -/
def readerRead (receiver : Option ReaderState) (destination : ReaderDestination) :
    Option ReaderCopyResult :=
  receiver.map fun state =>
    if (state.bytes.getD []).length ≤ state.cursor then
      ⟨0, .eof, destinationContents receiver destination, receiver⟩
    else
      let (count, output, updated) := copyReaderBytes state state.cursor destination
      ⟨count, .ok, output, some { updated with cursor := state.cursor + count, previousRune := -1 }⟩

/-- Go 1.27.0 src/bytes/reader.go:50–63. Signed offset validation precedes
    receiver access. Copy may change shared bytes, while reader fields stay fixed. -/
def readerReadAt (receiver : Option ReaderState) (offset : Int) (destination : ReaderDestination) :
    Option ReaderCopyResult :=
  if offset < 0 then some ⟨0, .negativeOffset, destinationContents receiver destination, receiver⟩
  else receiver.map fun state =>
    if (state.bytes.getD []).length ≤ offset.toNat then
      ⟨0, .eof, destinationContents receiver destination, receiver⟩
    else
      let (count, output, updated) := copyReaderBytes state offset.toNat destination
      ⟨count, if count < output.length then .eof else .ok, output, some updated⟩

theorem readerRead_at_end (s : ReaderState) (d : ReaderDestination)
    (atEnd : (s.bytes.getD []).length ≤ s.cursor) :
    readerRead (some s) d = some ⟨0, .eof, destinationContents (some s) d, some s⟩ := by
  simp [readerRead, atEnd]

theorem readerRead_shared :
    (readerRead (some ⟨some [1, 2, 3, 4], [9], 0, 7⟩) (.shared 1 4)).map
      (fun r => (r.count, r.destination, r.state.map (fun s => (s.bytes, s.spare, s.cursor, s.previousRune)))) =
    some (4, [1, 2, 3, 4], some (some [1, 1, 2, 3], [4], 4, -1)) := rfl

theorem readerRead_nil (d : ReaderDestination) : readerRead none d = none := rfl

theorem readerReadAt_negative (r : Option ReaderState) (offset : Int) (d : ReaderDestination)
    (negative : offset < 0) :
    readerReadAt r offset d = some ⟨0, .negativeOffset, destinationContents r d, r⟩ := by
  simp [readerReadAt, negative]

theorem readerReadAt_nil (offset : Int) (d : ReaderDestination) (nonnegative : 0 ≤ offset) :
    readerReadAt none offset d = none := by
  simp [readerReadAt, Int.not_lt.mpr nonnegative]

theorem readerReadAt_fields (s : ReaderState) (offset : Int) (d : ReaderDestination) :
    (readerReadAt (some s) offset d).bind
      (fun result => result.state.map (fun state => (state.cursor, state.previousRune))) =
      some (s.cursor, s.previousRune) := by
  by_cases negative : offset < 0
  · simp [readerReadAt, negative]
  · by_cases atEnd : (s.bytes.getD []).length ≤ offset.toNat
    · simp [readerReadAt, negative, atEnd]
    · cases d <;> simp [readerReadAt, negative, atEnd, copyReaderBytes]

theorem readerReadAt_partial :
    let r := readerReadAt (some ⟨some [1, 2, 3, 4], [9], 1, 7⟩) 2 (.shared 0 3)
    r.map (fun x => (x.count, x.status, x.destination,
      x.state.map (fun s => (s.bytes, s.spare, s.cursor, s.previousRune)))) =
      some (2, .eof, [3, 4, 3], some (some [3, 4, 3, 4], [9], 1, 7)) := rfl

end GoLean.Bytes
