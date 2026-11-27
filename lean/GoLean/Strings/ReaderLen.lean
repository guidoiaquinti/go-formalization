import Init.Data.BitVec.Basic
import Init.Data.List.Lemmas

namespace GoLean.Strings

/-- Public Reader operations keep the cursor nonnegative. The size of the
    backing string and the cursor fit in Go's respective int and int64.
    https://github.com/golang/go/blob/go1.27.0/src/strings/reader.go#L17-L30 -/
def readerLen (receiver : Option (List (BitVec 8) × Nat)) : Option Nat := do
  let (bytes, cursor) ← receiver
  return if cursor ≥ bytes.length then 0 else bytes.length - cursor

theorem readerLen_live (bytes : List (BitVec 8)) (cursor : Nat) :
    readerLen (some (bytes, cursor)) = some (bytes.drop cursor).length := by
  rw [List.length_drop]
  change some (if cursor ≥ bytes.length then 0 else bytes.length - cursor) = _
  apply congrArg some
  split <;> omega

theorem readerLen_nil : readerLen none = none := rfl

end GoLean.Strings
