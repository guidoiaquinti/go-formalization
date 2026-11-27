import GoLean.Strings.MakeStringFinder

namespace GoLean.Strings

private def firstMismatch (text pattern : Array (BitVec 8)) (right : Nat) :
    Nat → Option (Nat × Nat × BitVec 8)
  | 0 => none
  | j + 1 =>
    let i := right - 1
    let byte := text[i]!
    if byte = pattern[j]! then firstMismatch text pattern i j else some (i, j, byte)

private def searchFinder (finder : FinderTables) (text pattern : Array (BitVec 8))
    (right : Nat) : Nat → Int
  | 0 => -1
  | fuel + 1 =>
    if right > text.size then -1
    else
      match firstMismatch text pattern right pattern.size with
      | none => Int.ofNat (right - pattern.size)
      | some (i, j, byte) =>
        let shift := max finder.badCharSkip[byte.toNat]! finder.goodSuffixSkip[j]!
        searchFinder finder text pattern (i + shift + 1) fuel

/-- Go 1.27.0 src/strings/search.go:102–116 with tables initialized by
    makeStringFinder. The right cursor is one past the source's index before
    comparison, allowing natural indices even for the initial empty match. -/
def stringFinderNext (pattern text : List (BitVec 8)) : Int :=
  let tables := makeStringFinder pattern
  searchFinder tables text.toArray pattern.toArray pattern.length (text.length + 1)

theorem stringFinderNext_empty (text : List (BitVec 8)) :
    stringFinderNext [] text = 0 := by
  simp [stringFinderNext, searchFinder, firstMismatch]

theorem stringFinderNext_tooShort (pattern text : List (BitVec 8))
    (short : text.length < pattern.length) :
    stringFinderNext pattern text = -1 := by
  simp [stringFinderNext, searchFinder, short]

private theorem firstMismatch_self (bytes : Array (BitVec 8)) (right : Nat) :
    firstMismatch bytes bytes right right = none := by
  induction right with
  | zero => rfl
  | succ right ih => simpa [firstMismatch] using ih

theorem stringFinderNext_self (pattern : List (BitVec 8)) :
    stringFinderNext pattern pattern = 0 := by
  simp [stringFinderNext, searchFinder, firstMismatch_self]

/-- Search with a previously constructed finder instead of rebuilding its tables. -/
def FinderTables.next (finder : FinderTables) (text : List (BitVec 8)) : Int :=
  searchFinder finder text.toArray finder.pattern.toArray finder.pattern.length (text.length + 1)

theorem FinderTables.next_constructed (pattern text : List (BitVec 8)) :
    (makeStringFinder pattern).next text = stringFinderNext pattern text := rfl

end GoLean.Strings
