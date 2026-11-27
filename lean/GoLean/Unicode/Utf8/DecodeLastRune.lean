import GoLean.Unicode.Utf8.DecodeLastRuneInString
import GoLean.Unicode.Utf8.DecodeRune

namespace GoLean.Unicode.Utf8

/-- Go 1.27.0 src/unicode/utf8/utf8.go:279–306. The slice decoder
    scans the same positions as the string decoder, then calls DecodeRune. -/
def decodeLastRune (bytes : List Nat) : Nat × Nat :=
  if bytes.isEmpty then (0xfffd, 0)
  else
    let last := bytes[bytes.length - 1]!
    if last < 128 then (last, 1)
    else
      let limit := bytes.length - 4
      let start := lastRuneStart bytes limit (bytes.length - 1 - limit)
      let decoded := decodeRune (bytes.drop start)
      if start + decoded.2 = bytes.length then decoded else (0xfffd, 1)

theorem decodeLastRune_spec (bytes : List Nat) :
    decodeLastRune bytes = decodeLastRuneInString bytes := by
  simp only [decodeLastRune, decodeLastRuneInString, decodeRune_spec]
  split
  · rfl
  · split <;> rfl

end GoLean.Unicode.Utf8
