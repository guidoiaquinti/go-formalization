import GoLean.Compress.Flate.DictionaryState

namespace GoLean.Compress.Flate

structure DictionaryInit where
  next : DictionaryState
  oldBacking : List (BitVec 8)
  reused : Bool

/-- Retain the suffix of an oversized preset and leave spare reused bytes intact. -/
def dictionaryInit (s : DictionaryState) (size : Nat) (preset : List (BitVec 8)) :
    DictionaryInit :=
  let reused := size ≤ s.backing.length
  let backing := if reused then s.backing else List.replicate size 0
  let used := preset.drop (preset.length - size)
  let output := used ++ backing.drop used.length
  let full := used.length == size
  let cursor := if full then 0 else used.length
  ⟨⟨output, size, cursor, cursor, full, !reused || s.present⟩,
    if reused then output else s.backing, reused⟩

theorem dictionaryInit_zero (s : DictionaryState) (preset : List (BitVec 8)) :
    (dictionaryInit s 0 preset).next =
      { s with length := 0, read := 0, written := 0, full := true } := by
  simp [dictionaryInit]

theorem dictionaryInit_readEqualsWrite (s : DictionaryState) (size : Nat)
    (preset : List (BitVec 8)) :
    (dictionaryInit s size preset).next.read = (dictionaryInit s size preset).next.written := rfl

end GoLean.Compress.Flate
