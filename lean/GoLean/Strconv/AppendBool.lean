import GoLean.Strconv.FormatBool

namespace GoLean.Strconv

/-- Go 1.27.0 src/strconv/number.go:32-36 delegates to
src/internal/strconv/atob.go:28-35. Lists represent the visible slice bytes. -/
def appendBool (dst : List UInt8) (b : Bool) : List UInt8 :=
  if b then dst ++ "true".toUTF8.data.toList
  else dst ++ "false".toUTF8.data.toList

theorem appendBool_spec (dst : List UInt8) (b : Bool) :
    appendBool dst b = dst ++ (formatBool b).toUTF8.data.toList := by
  cases b <;> rfl


end GoLean.Strconv
