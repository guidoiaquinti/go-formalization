namespace GoLean.Strconv

structure BoolSyntaxError where
  func : String
  num : List UInt8

private def trueSpellings := ["1", "t", "T", "true", "TRUE", "True"]
private def falseSpellings := ["0", "f", "F", "false", "FALSE", "False"]
private def spellingBytes (s : String) := s.toUTF8.data.toList

/-- Go 1.27.0 src/strconv/number.go:19-25 wraps the syntax error from
src/internal/strconv/atob.go:10-18 through toError at number.go:231-235. The optional payload denotes a NumError
with ErrSyntax. The raw bytes also represent non-UTF-8 Go strings. -/
def parseBool (input : List UInt8) : Bool × Option BoolSyntaxError :=
  if (trueSpellings.map spellingBytes).contains input then (true, none)
  else if (falseSpellings.map spellingBytes).contains input then (false, none)
  else (false, some ⟨"ParseBool", input⟩)

theorem parseBool_spec (input : List UInt8) :
    ((parseBool input).1 = true ↔ input ∈ trueSpellings.map spellingBytes) ∧
    ((parseBool input).2 = none ↔
      input ∈ trueSpellings.map spellingBytes ∨ input ∈ falseSpellings.map spellingBytes) ∧
    ((parseBool input).2 ≠ none → parseBool input = (false, some ⟨"ParseBool", input⟩)) := by
  by_cases ht : input ∈ trueSpellings.map spellingBytes <;>
    by_cases hf : input ∈ falseSpellings.map spellingBytes <;>
    simp [parseBool, ht, hf]

end GoLean.Strconv
