import GoLean.Strconv.AppendBool
import GoLean.Strconv.ParseBool

namespace GoLean.Internal.Strconv

/-- Go 1.27.0 src/internal/strconv/atob.go returns bare ErrSyntax. -/
def parseBool (raw : List UInt8) : Bool × Bool :=
  let (value, error) := GoLean.Strconv.parseBool raw
  (value, error.isSome)

def formatBool := GoLean.Strconv.formatBool

structure BoolAppend where
  result : List UInt8
  oldBacking : List UInt8
  reused : Bool

/-- The original backing includes spare capacity; allocation size is abstracted. -/
def appendBool (backing : List UInt8) (length : Nat) (value : Bool) : BoolAppend :=
  let result := GoLean.Strconv.appendBool (backing.take length) value
  if result.length ≤ backing.length then
    ⟨result, result ++ backing.drop result.length, true⟩
  else ⟨result, backing, false⟩

/-- Go 1.27.0 src/internal/strconv/atoi.go:25 for every signed Error tag. -/
def errorText : Int → String
  | 1 => "value out of range"
  | 2 => "invalid syntax"
  | 3 => "invalid base"
  | 4 => "invalid bit size"
  | _ => "unknown error"

theorem parseBool_false : parseBool [70, 97, 108, 115, 101] = (false, false) := by decide

theorem parseBool_invalid : parseBool [84, 114, 85, 101] = (false, true) := by decide

theorem formatBool_spellings :
    formatBool false = "false" ∧ formatBool true = "true" := ⟨rfl, rfl⟩

theorem appendBool_short :
    (appendBool [19, 20, 21, 22, 23] 1 false).oldBacking = [19, 20, 21, 22, 23] ∧
    (appendBool [19, 20, 21, 22, 23] 1 false).reused = false := by decide

theorem appendBool_exact :
    (appendBool [19, 20, 21, 22, 23] 1 true).oldBacking = [19, 116, 114, 117, 101] ∧
    (appendBool [19, 20, 21, 22, 23] 1 true).reused = true := by decide

theorem errorText_syntax : errorText 2 = "invalid syntax" := rfl

end GoLean.Internal.Strconv
