namespace GoLean.Unicode.Utf8

/-- A Go rune is interpreted as a signed integer. The two branches use the
    scalar bounds and surrogate constants from the pinned source.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf8/utf8.go#L14-L26
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf8/utf8.go#L568-L578 -/
def validRune (r : Int) : Bool :=
  if 0 ≤ r ∧ r < 0xd800 then true
  else if 0xdfff < r ∧ r ≤ 0x10ffff then true
  else false

/-- A legal scalar is on a Unicode plane and has neither surrogate prefix. -/
theorem validRune_spec (r : Int) :
    validRune r = true ↔
      (0 ≤ r / 65536 ∧ r / 65536 ≤ 16) ∧ r / 1024 ≠ 54 ∧ r / 1024 ≠ 55 := by
  unfold validRune
  split
  · simp <;> omega
  · split <;> simp <;> omega

end GoLean.Unicode.Utf8
