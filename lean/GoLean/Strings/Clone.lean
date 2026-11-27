import Init.Data.BitVec.Basic

namespace GoLean.Strings

/-- The empty branch returns the canonical empty value. Otherwise Go creates
    a buffer, copies the bytes, and returns a string over that buffer.
    The Boolean records the fresh backing-storage branch.
    https://github.com/golang/go/blob/go1.27.0/src/strings/clone.go#L11-L23
    https://github.com/golang/go/blob/go1.27.0/src/internal/stringslite/strings.go#L115-L122 -/
def clone (bytes : List (BitVec 8)) : List (BitVec 8) × Bool :=
  if bytes.isEmpty then ([], false) else (bytes.map id, true)

theorem clone_spec (bytes : List (BitVec 8)) :
    clone bytes = (bytes, decide (0 < bytes.length)) := by
  cases bytes <;> simp [clone]

end GoLean.Strings
