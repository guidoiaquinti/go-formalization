import GoLean.Bytes.Count
import GoLean.Strings.Replace

namespace GoLean.Bytes

structure ByteCopy where
  bytes : Option (List Nat)
  exactCapacity : Option Nat

/-- Go 1.27.0 src/bytes/bytes.go:1391–1396. Nil remains nil. The
    append result for populated input has allocator-dependent capacity. -/
def cloneBytes (input : Option (List Nat)) : ByteCopy :=
  ⟨input, if (input.getD []).isEmpty then some 0 else none⟩

/-- Go 1.27.0 src/bytes/bytes.go:1162–1200. A no-match append uses a
    nil base. Successful replacement starts with a nonnil exact allocation,
    including when all bytes are removed. The raw-byte transformation is
    shared with the string model, which also splits empty-pattern matches
    at UTF-8 rune boundaries. -/
def replaceBytes (input : Option (List Nat)) (old replacement : List Nat) (limit : Int) : ByteCopy :=
  let bytes := input.getD []
  let found := if limit == 0 then 0 else count bytes old
  if found == 0 then
    if bytes.isEmpty then ⟨none, some 0⟩ else ⟨some bytes, none⟩
  else
    let n := if limit < 0 then found else min found limit.toNat
    let capacity := Int.ofNat bytes.length + Int.ofNat n *
      (Int.ofNat replacement.length - Int.ofNat old.length)
    ⟨some (GoLean.Strings.replace bytes old replacement limit), some capacity.toNat⟩

/-- Go 1.27.0 src/bytes/bytes.go:1206–1208 delegates with a negative limit. -/
def replaceAllBytes (input : Option (List Nat)) (old replacement : List Nat) : ByteCopy :=
  replaceBytes input old replacement (-1)

theorem cloneBytes_contents (input : Option (List Nat)) :
    (cloneBytes input).bytes = input := rfl

theorem cloneBytes_empty : cloneBytes (some []) = ⟨some [], some 0⟩ := rfl

theorem replaceBytes_zero (input : Option (List Nat)) (old replacement : List Nat) :
    replaceBytes input old replacement 0 =
      let bytes := input.getD []
      if bytes.isEmpty then ⟨none, some 0⟩ else ⟨some bytes, none⟩ := by
  simp [replaceBytes]

theorem replaceBytes_removed_all :
    replaceBytes (some [97, 97]) [97] [] (-1) = ⟨some [], some 0⟩ := rfl

theorem replaceBytes_emptyPattern :
    replaceBytes none [] [] 1 = ⟨some [], some 0⟩ := rfl

theorem replaceAllBytes_unlimited (input : Option (List Nat)) (old replacement : List Nat) :
    replaceAllBytes input old replacement = replaceBytes input old replacement (-1) := rfl

theorem replaceAllBytes_malformed :
    replaceAllBytes (some [0xe2, 0x82, 0xef, 0xbf, 0xbd]) [] [0xff] =
      ⟨some [0xff, 0xe2, 0xff, 0x82, 0xff, 0xef, 0xbf, 0xbd, 0xff], some 9⟩ := by
  have matching : count [0xe2, 0x82, 0xef, 0xbf, 0xbd] [] = 4 := by
    rw [count_empty]
    have first : GoLean.Strings.runeWidth 226 [130, 239, 191, 189] = 1 := rfl
    have second : GoLean.Strings.runeWidth 130 [239, 191, 189] = 1 := rfl
    have third : GoLean.Strings.runeWidth 239 [191, 189] = 3 := rfl
    simp [GoLean.Strings.countRunes, first, second, third]
  simp [replaceAllBytes, replaceBytes, matching, GoLean.Strings.replace_emptyMalformed]

end GoLean.Bytes
