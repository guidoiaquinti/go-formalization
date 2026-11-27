import GoLean.Unicode.Utf8.AppendRune

namespace GoLean.Unicode.Utf8

/-- Go 1.27.0 src/unicode/utf8/utf8.go:365–405. ASCII uses its indexed
    assignment, and each longer encoding checks the final destination
    index before any write. None denotes the bounds panic. -/
def encodeRune (dst : List Nat) (r : Int) : Option Nat × List Nat :=
  let bytes := appendRune [] r
  if bytes.length ≤ dst.length then
    (some bytes.length, bytes ++ dst.drop bytes.length)
  else
    (none, dst)

theorem encodeRune_short (dst : List Nat) (r : Int)
    (short : dst.length < (appendRune [] r).length) :
    encodeRune dst r = (none, dst) := by
  simp [encodeRune, Nat.not_le_of_lt short]

theorem encodeRune_written (dst : List Nat) (r : Int)
    (enough : (appendRune [] r).length ≤ dst.length) :
    let bytes := appendRune [] r
    let result := encodeRune dst r
    result.1 = some bytes.length ∧
      result.2.take bytes.length = bytes ∧
      result.2.drop bytes.length = dst.drop bytes.length := by
  simp [encodeRune, enough]

end GoLean.Unicode.Utf8
