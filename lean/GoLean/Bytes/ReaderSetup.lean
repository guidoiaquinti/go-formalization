import Init.Data.BitVec.Basic

namespace GoLean.Bytes

structure ReaderState where
  bytes : Option (List (BitVec 8))
  spare : List (BitVec 8)
  cursor : Nat
  previousRune : Int

/-- Go 1.27.0 NewReader retains the input slice and initializes the cursor
    and previous-rune marker. -/
def newReader (bytes : Option (List (BitVec 8))) (spare : List (BitVec 8)) : ReaderState :=
  ⟨bytes, spare, 0, -1⟩

def readerReset (receiver : Option ReaderState) (bytes : Option (List (BitVec 8)))
    (spare : List (BitVec 8)) : Option ReaderState :=
  receiver.map fun _ => newReader bytes spare

def readerSize (receiver : Option ReaderState) : Option Nat :=
  receiver.map fun state => (state.bytes.getD []).length

theorem newReader_snapshot (bytes : Option (List (BitVec 8))) (spare : List (BitVec 8)) :
    newReader bytes spare = ⟨bytes, spare, 0, -1⟩ := rfl

theorem readerReset_snapshot (old : ReaderState) (bytes : Option (List (BitVec 8)))
    (spare : List (BitVec 8)) :
    readerReset (some old) bytes spare = some ⟨bytes, spare, 0, -1⟩ := rfl

theorem readerReset_nil (bytes : Option (List (BitVec 8))) (spare : List (BitVec 8)) :
    readerReset none bytes spare = none := rfl

theorem readerSize_unreadIndependent (bytes : Option (List (BitVec 8))) (spare : List (BitVec 8))
    (cursor : Nat) (previousRune : Int) :
    readerSize (some ⟨bytes, spare, cursor, previousRune⟩) = some (bytes.getD []).length := rfl

theorem readerSize_nil : readerSize none = none := rfl

end GoLean.Bytes
