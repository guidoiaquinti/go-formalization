namespace GoLean.Strconv

/-- Go 1.27.0 src/internal/strconv/atob.go:21-26. -/
private def formatBoolImpl (b : Bool) : String :=
  if b then "true" else "false"

/-- Go 1.27.0 src/strconv/number.go:28-30 delegates to internal/strconv. -/
def formatBool (b : Bool) : String := formatBoolImpl b

theorem formatBool_spec (b : Bool) :
    formatBool b = (match b with | false => "false" | true => "true") := by
  cases b <;> rfl


end GoLean.Strconv
