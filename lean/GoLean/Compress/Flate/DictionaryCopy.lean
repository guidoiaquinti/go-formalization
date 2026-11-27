import GoLean.Compress.Flate.DictionaryState

namespace GoLean.Compress.Flate

private def copyForward (size distance : Nat) :
    Nat → Nat → List (BitVec 8) → List (BitVec 8)
  | 0, _, backing => backing
  | count + 1, position, backing =>
    let source := (position + size - distance) % size
    copyForward size distance count (position + 1) (backing.set position backing[source]!)

private def doCopy (state : DictionaryState) (dist count : Nat) : Nat × DictionaryState :=
  (count, { state with
    backing := copyForward state.length dist count state.written state.backing
    written := state.written + count })

/-- Legal positive distances refer to history, including history across the ring boundary. -/
def writeCopy (state : DictionaryState) (dist count : Nat) : Nat × DictionaryState :=
  doCopy state dist (min count (availWrite state))

def tryWriteCopy (state : DictionaryState) (dist count : Nat) : Nat × DictionaryState :=
  if state.written < dist || availWrite state < count then (0, state)
  else doCopy state dist count

theorem writeCopy_space (s : DictionaryState) (dist count : Nat) :
    (writeCopy s dist count).1 = min count (availWrite s) := rfl

theorem tryWriteCopy_wrapped (s : DictionaryState) (dist count : Nat) (wrapped : s.written < dist) :
    tryWriteCopy s dist count = (0, s) := by simp [tryWriteCopy, wrapped]

theorem tryWriteCopy_short (s : DictionaryState) (dist count : Nat) (short : availWrite s < count) :
    tryWriteCopy s dist count = (0, s) := by simp [tryWriteCopy, short]

end GoLean.Compress.Flate
