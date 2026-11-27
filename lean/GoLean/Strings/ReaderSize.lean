import Init.Data.BitVec.Basic

namespace GoLean.Strings

/-- The backing byte length fits in both Go's int and int64. The Reader
    cursor is present in the receiver state but does not affect the result.
    https://github.com/golang/go/blob/go1.27.0/src/strings/reader.go#L17-L36 -/
def readerSize (receiver : Option (List (BitVec 8) × Nat)) : Option Nat := do
  let (bytes, _) ← receiver
  return bytes.length

theorem readerSize_live (bytes : List (BitVec 8)) (cursor : Nat) :
    readerSize (some (bytes, cursor)) = some bytes.length := rfl

theorem readerSize_nil : readerSize none = none := rfl

end GoLean.Strings
