import Init.Data.BitVec.Lemmas

namespace GoLean.Unicode.Utf8

/-- Go 1.27.0, src/unicode/utf8/utf8.go, lines 454–457.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf8/utf8.go#L454-L457 -/
def runeStart (b : BitVec 8) : Bool :=
  (b &&& 0xC0) != 0x80

/-- Only bytes in the UTF-8 continuation interval are rejected. -/
theorem runeStart_spec : ∀ b : BitVec 8,
    runeStart b = (decide (b < 128) || decide (192 ≤ b)) := by
  decide

end GoLean.Unicode.Utf8
