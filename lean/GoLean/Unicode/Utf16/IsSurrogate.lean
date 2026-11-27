namespace GoLean.Unicode.Utf16

/-- Go 1.27.0, src/unicode/utf16/utf16.go, lines 17–32.
    A Go rune is interpreted as its signed mathematical integer.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf16/utf16.go#L17-L32 -/
def isSurrogate (r : Int) : Bool :=
  decide (0xd800 ≤ r) && decide (r < 0xe000)

theorem isSurrogate_spec (r : Int) :
    isSurrogate r = true ↔ r / 1024 = 54 ∨ r / 1024 = 55 := by
  simp only [isSurrogate, Bool.and_eq_true, decide_eq_true_eq]
  omega

end GoLean.Unicode.Utf16
