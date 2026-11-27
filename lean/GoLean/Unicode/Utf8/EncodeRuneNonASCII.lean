import GoLean.Unicode.Utf8.AppendRuneNonASCII

namespace GoLean.Unicode.Utf8

/-- Go 1.27.0 src/unicode/utf8/utf8.go:377–405. Every path checks the final
    destination index before any write. None denotes the bounds panic. -/
def encodeRuneNonASCII (dst : List Nat) (r : Int) : Option Nat × List Nat :=
  let bytes := appendRuneNonASCII [] r
  if bytes.length ≤ dst.length then
    (some bytes.length, bytes ++ dst.drop bytes.length)
  else
    (none, dst)

theorem encodeRuneNonASCII_short (dst : List Nat) (r : Int)
    (short : dst.length < (appendRuneNonASCII [] r).length) :
    encodeRuneNonASCII dst r = (none, dst) := by
  simp [encodeRuneNonASCII, Nat.not_le_of_lt short]

theorem encodeRuneNonASCII_written (dst : List Nat) (r : Int)
    (enough : (appendRuneNonASCII [] r).length ≤ dst.length) :
    let bytes := appendRuneNonASCII [] r
    let result := encodeRuneNonASCII dst r
    result.1 = some bytes.length ∧
      result.2.take bytes.length = bytes ∧
      result.2.drop bytes.length = dst.drop bytes.length := by
  simp [encodeRuneNonASCII, enough]

end GoLean.Unicode.Utf8
