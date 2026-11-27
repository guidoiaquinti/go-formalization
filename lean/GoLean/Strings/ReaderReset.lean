import Init.Data.BitVec.Basic

namespace GoLean.Strings

structure ReaderState where
  contents : List (BitVec 8)
  cursor : Nat
  previousRune : Int

/-- Reset replaces all fields of a live Reader. A nil receiver panics.
    https://github.com/golang/go/blob/go1.27.0/src/strings/reader.go#L17-L22
    https://github.com/golang/go/blob/go1.27.0/src/strings/reader.go#L151-L152 -/
def readerReset (receiver : Option ReaderState) (contents : List (BitVec 8)) :
    Option ReaderState := do
  let _ ← receiver
  return ⟨contents, 0, -1⟩

theorem readerReset_live (old : ReaderState) (contents : List (BitVec 8)) :
    readerReset (some old) contents = some ⟨contents, 0, -1⟩ := rfl

theorem readerReset_nil (contents : List (BitVec 8)) :
    readerReset none contents = none := rfl

end GoLean.Strings
