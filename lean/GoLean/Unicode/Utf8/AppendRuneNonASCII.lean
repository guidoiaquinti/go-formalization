import GoLean.Unicode.Utf8.AppendRune

namespace GoLean.Unicode.Utf8

/-- Go 1.27.0 src/unicode/utf8/utf8.go:418–430. A direct ASCII call
    enters the two-byte branch and deliberately differs from AppendRune. -/
def appendRuneNonASCII (dst : List Nat) (r : Int) : List Nat :=
  if 0 ≤ r ∧ r < 128 then dst ++ [192 + r.toNat / 64, 128 + r.toNat % 64]
  else appendRune dst r

theorem appendRuneNonASCII_ascii (dst : List Nat) (r : Int) (ascii : 0 ≤ r ∧ r < 128) :
    appendRuneNonASCII dst r = dst ++ [192 + r.toNat / 64, 128 + r.toNat % 64] := by
  simp [appendRuneNonASCII, ascii]

theorem appendRuneNonASCII_non_ascii (dst : List Nat) (r : Int)
    (nonASCII : ¬ (0 ≤ r ∧ r < 128)) :
    (appendRuneNonASCII dst r).take dst.length = dst ∧
      (appendRuneNonASCII dst r).drop dst.length =
        Strings.encodeRune (if validRune r then r.toNat else 0xfffd) := by
  simpa [appendRuneNonASCII, nonASCII] using appendRune_spec dst r

end GoLean.Unicode.Utf8
