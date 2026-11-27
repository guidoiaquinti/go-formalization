import GoLean.Strings.IndexRune
import GoLean.Unicode.Utf8.ValidRune

namespace GoLean.Unicode.Utf8

/-- The source appends the encoded scalar bytes, replacing an invalid rune.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf8/utf8.go#L407-L430 -/
def appendRune (dst : List Nat) (r : Int) : List Nat :=
  dst ++ if validRune r then Strings.encodeRune r.toNat else [0xef, 0xbf, 0xbd]

theorem appendRune_spec (dst : List Nat) (r : Int) :
    (appendRune dst r).take dst.length = dst ∧
    (appendRune dst r).drop dst.length =
      Strings.encodeRune (if validRune r then r.toNat else 0xfffd) := by
  constructor
  · simp [appendRune]
  · by_cases legal : validRune r = true
    · simp [appendRune, legal]
    · have invalid : validRune r = false := Bool.eq_false_iff.mpr legal
      simp only [appendRune, invalid, Bool.false_eq_true, ↓reduceIte, List.drop_left]
      rfl

end GoLean.Unicode.Utf8
