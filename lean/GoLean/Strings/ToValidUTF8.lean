import GoLean.Bytes.ToValidUTF8

namespace GoLean.Strings

/-- Go 1.27.0 src/strings/strings.go:771–817. Project the shared scanner's
    bytes. A valid prefix is copied unchanged whether Go returns the original
    string or initializes a builder at the first malformed byte. -/
def toValidUTF8 (input replacement : List Nat) : List Nat :=
  (GoLean.Bytes.toValidUTF8 input replacement).1

theorem toValidUTF8_ascii (input replacement : List Nat) (ascii : ∀ byte ∈ input, byte < 128) :
    toValidUTF8 input replacement = input := by
  induction input with
  | nil => simp [toValidUTF8, GoLean.Bytes.toValidUTF8_empty]
  | cons head tail ih =>
    have headASCII : head < 128 := ascii head (by simp)
    have tailASCII : ∀ byte ∈ tail, byte < 128 := by
      intro byte h
      exact ascii byte (by simp [h])
    unfold toValidUTF8 at ih ⊢
    rw [GoLean.Bytes.toValidUTF8_ascii _ _ _ headASCII, ih tailASCII]

theorem toValidUTF8_invalidRun (replacement : List Nat) :
    toValidUTF8 [0xff, 0x80, 0xc0] replacement = replacement := by
  change replacement ++ [] = replacement
  simp

theorem toValidUTF8_encodedRuneError (replacement : List Nat) :
    toValidUTF8 [0xff, 0xef, 0xbf, 0xbd, 0x80] replacement =
      replacement ++ [0xef, 0xbf, 0xbd] ++ replacement := by
  change replacement ++ ([0xef, 0xbf, 0xbd] ++ (replacement ++ [])) = _
  simp [List.append_assoc]

end GoLean.Strings
