namespace GoLean.Unicode.Utf16

/-- The signed rune is represented by its mathematical integer value. The
    conditions occur in the same order as the cases in the Go implementation.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf16/utf16.go#L10-L26
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf16/utf16.go#L55-L66 -/
def runeLen (r : Int) : Int :=
  if (0 ≤ r ∧ r < 0xd800) ∨ (0xe000 ≤ r ∧ r < 0x10000) then 1
  else if 0x10000 ≤ r ∧ r ≤ 0x10ffff then 2
  else -1

/-- Unicode has 17 planes. A valid plane-zero rune needs one word, and a
    valid rune in any other plane needs two. Blocks 54 and 55 are reserved. -/
theorem runeLen_spec (r : Int) :
    runeLen r =
      if r < 0 ∨ 17 ≤ r / 65536 ∨ r / 1024 = 54 ∨ r / 1024 = 55 then -1
      else if r / 65536 = 0 then 1 else 2 := by
  unfold runeLen
  split <;> split
  all_goals try split
  all_goals omega

end GoLean.Unicode.Utf16
