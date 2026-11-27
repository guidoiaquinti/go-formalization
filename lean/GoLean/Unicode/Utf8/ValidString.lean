import GoLean.Unicode.Utf8.Valid

namespace GoLean.Unicode.Utf8

/-- Go 1.27.0 src/unicode/utf8/utf8.go:521–567. Validation uses the
    same native-word ASCII checks and multibyte paths as the slice validator.
    A Go string enters the shared model as its raw bytes. -/
def validString (wordWidth : Nat) (bytes : List Nat) : Bool :=
  valid wordWidth bytes

theorem validString_ascii_only (wordWidth : Nat) (bytes : List Nat)
    (ascii : ∀ b ∈ bytes, b < 128) : validString wordWidth bytes = true :=
  valid_ascii_only wordWidth bytes ascii

theorem validString_rejects_malformed_lead (wordWidth lead : Nat) (rest : List Nat)
    (nonASCII : 128 ≤ lead) (malformed : GoLean.Strings.runeWidth lead rest = 1) :
    validString wordWidth (lead :: rest) = false :=
  valid_rejects_malformed_lead wordWidth lead rest nonASCII malformed

end GoLean.Unicode.Utf8
