namespace GoLean.Unicode.Utf8

private def packByteLanes : List (BitVec 8) → Nat
  | [] => 0
  | byte :: tail => byte.toNat + 256 * packByteLanes tail

/-- Go 1.27.0 src/unicode/utf8/utf8.go:462–467. The bitwise OR of
    disjoint byte lanes is represented by addition. None is an index panic. -/
def word (width : Nat) (bytes : List (BitVec 8)) : Option Nat :=
  if width ≤ bytes.length then some (packByteLanes (bytes.take width)) else none

theorem word_short (width : Nat) (bytes : List (BitVec 8)) (short : bytes.length < width) :
    word width bytes = none := by simp [word, Nat.not_le_of_lt short]

theorem word_suffix_ignored (width : Nat) (front suffix : List (BitVec 8))
    (enough : width ≤ front.length) : word width (front ++ suffix) = word width front := by
  have extended : width ≤ front.length + suffix.length := by omega
  simp [word, enough, extended, List.take_append_of_le_length enough]

theorem word_first_lane (width : Nat) (byte : BitVec 8) (tail : List (BitVec 8))
    (enough : width ≤ tail.length) :
    (word (width + 1) (byte :: tail)).map (· % 256) = some byte.toNat := by
  have extended : width + 1 ≤ (byte :: tail).length := by simp [enough]
  simp only [word, extended, ↓reduceIte, List.take_succ_cons, packByteLanes, Option.map_some]
  have small := byte.isLt
  congr 1
  omega

end GoLean.Unicode.Utf8
