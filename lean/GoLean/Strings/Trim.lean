import GoLean.Strings.MakeASCIISet
import GoLean.Strings.IndexRune
import GoLean.Unicode.Utf8.DecodeLastRuneInString

namespace GoLean.Strings

private def trimMatchingBytes (input : List (BitVec 8)) (member : BitVec 8 → Bool) :
    List (BitVec 8) :=
  ((input.reverse.dropWhile member).reverse).dropWhile member

private def scanTrimLeftRunes (member : Nat → Bool) (input : List Nat) : Nat → List Nat
  | 0 => input
  | fuel + 1 =>
    if input.isEmpty then input
    else
      let (rune, width) := GoLean.Unicode.Utf8.decodeRuneInString input
      if member rune then scanTrimLeftRunes member (input.drop width) fuel else input

private def scanTrimRightRunes (member : Nat → Bool) (input : List Nat) : Nat → List Nat
  | 0 => input
  | fuel + 1 =>
    match input.getLast? with
    | none => input
    | some last =>
      let (rune, width) := if last < 128 then (last, 1)
        else GoLean.Unicode.Utf8.decodeLastRuneInString input
      if member rune then scanTrimRightRunes member (input.take (input.length - width)) fuel
      else input

private def trimBySet (input cutset : List (BitVec 8)) : List (BitVec 8) :=
  let (table, ascii) := makeASCIISet cutset
  if ascii then trimMatchingBytes input table
  else
    let member := fun rune => containsRune (cutset.map BitVec.toNat) (Int.ofNat rune)
    let raw := input.map BitVec.toNat
    let remaining := scanTrimRightRunes member raw raw.length
    (scanTrimLeftRunes member remaining remaining.length).map (BitVec.ofNat 8)

/-- Go 1.27.0 src/strings/strings.go:980–991, 1010–1036, 1055–1085.
    Trim the right edge before the left. Single ASCII cutsets use byte
    equality, larger ASCII cutsets use the lookup set, and other cutsets
    decode UTF-8 and test rune membership without rewriting retained bytes. -/
def trim (input cutset : List (BitVec 8)) : List (BitVec 8) :=
  if input.isEmpty || cutset.isEmpty then input
  else
    match cutset with
    | [byte] => if byte.toNat < 128 then trimMatchingBytes input (· == byte) else trimBySet input cutset
    | _ => trimBySet input cutset

private theorem dropWhileLen (p : α → Bool) (input : List α) :
    (input.dropWhile p).length ≤ input.length := by
  induction input with
  | nil => simp
  | cons x xs ih =>
    simp only [List.dropWhile, List.length_cons]
    split
    · exact Nat.le_trans ih (Nat.le_succ _)
    · simp

private theorem trimMatchingBytesLen (input : List (BitVec 8)) (member : BitVec 8 → Bool) :
    (trimMatchingBytes input member).length ≤ input.length := by
  unfold trimMatchingBytes
  have right := dropWhileLen member input.reverse
  have left := dropWhileLen member (input.reverse.dropWhile member).reverse
  simp only [List.length_reverse] at right left
  exact Nat.le_trans left right

private theorem scanTrimLeftRunesLen (member : Nat → Bool) (input : List Nat) (fuel : Nat) :
    (scanTrimLeftRunes member input fuel).length ≤ input.length := by
  induction fuel generalizing input with
  | zero => simp [scanTrimLeftRunes]
  | succ fuel ih =>
    unfold scanTrimLeftRunes
    split
    · exact Nat.le_refl _
    · dsimp only
      split
      · exact Nat.le_trans (ih _) (by simp)
      · exact Nat.le_refl _

private theorem scanTrimRightRunesLen (member : Nat → Bool) (input : List Nat) (fuel : Nat) :
    (scanTrimRightRunes member input fuel).length ≤ input.length := by
  induction fuel generalizing input with
  | zero => simp [scanTrimRightRunes]
  | succ fuel ih =>
    unfold scanTrimRightRunes
    split
    · exact Nat.le_refl _
    · dsimp only
      split <;> split
      all_goals first | exact Nat.le_trans (ih _) (List.length_take_le' _ input) | exact Nat.le_refl _

private theorem trimBySetLen (input cutset : List (BitVec 8)) :
    (trimBySet input cutset).length ≤ input.length := by
  unfold trimBySet
  dsimp only
  split
  · exact trimMatchingBytesLen _ _
  · simp only [List.length_map]
    exact Nat.le_trans (scanTrimLeftRunesLen _ _ _) (by
      simpa using scanTrimRightRunesLen _ (input.map BitVec.toNat) _)

theorem trim_length (input cutset : List (BitVec 8)) :
    (trim input cutset).length ≤ input.length := by
  unfold trim
  split
  · exact Nat.le_refl _
  · split
    · split
      · exact trimMatchingBytesLen _ _
      · exact trimBySetLen _ _
    · exact trimBySetLen _ _

theorem trim_emptyCutset (input : List (BitVec 8)) :
    trim input [] = input := by
  simp [trim]

theorem trim_emptyInput (cutset : List (BitVec 8)) :
    trim [] cutset = [] := by
  simp [trim]

theorem trim_singleASCII (input : List (BitVec 8)) (byte : BitVec 8)
    (ascii : byte.toNat < 128) :
    trim input [byte] =
      ((input.reverse.dropWhile (· == byte)).reverse).dropWhile (· == byte) := by
  cases input <;> simp [trim, trimMatchingBytes, ascii]

private def trimLeftBySet (input cutset : List (BitVec 8)) : List (BitVec 8) :=
  let (table, ascii) := makeASCIISet cutset
  if ascii then input.dropWhile table
  else
    let member := fun rune => containsRune (cutset.map BitVec.toNat) (Int.ofNat rune)
    let raw := input.map BitVec.toNat
    (scanTrimLeftRunes member raw raw.length).map (BitVec.ofNat 8)

/-- Go 1.27.0 src/strings/strings.go:997–1036. Dispatch by cutset
    representation and remove leading members while retaining raw bytes. -/
def trimLeft (input cutset : List (BitVec 8)) : List (BitVec 8) :=
  if input.isEmpty || cutset.isEmpty then input
  else
    match cutset with
    | [byte] => if byte.toNat < 128 then input.dropWhile (· == byte) else trimLeftBySet input cutset
    | _ => trimLeftBySet input cutset

theorem trimLeft_emptyCutset (input : List (BitVec 8)) :
    trimLeft input [] = input := by
  simp [trimLeft]

theorem trimLeft_singleASCII (input : List (BitVec 8)) (byte : BitVec 8)
    (ascii : byte.toNat < 128) :
    trimLeft input [byte] = input.dropWhile (· == byte) := by
  cases input <;> simp [trimLeft, ascii]

theorem trimLeft_length (input cutset : List (BitVec 8)) :
    (trimLeft input cutset).length ≤ input.length := by
  have setLength : (trimLeftBySet input cutset).length ≤ input.length := by
    unfold trimLeftBySet
    dsimp only
    split
    · exact dropWhileLen _ _
    · simpa using scanTrimLeftRunesLen _ (input.map BitVec.toNat) _
  unfold trimLeft
  split
  · exact Nat.le_refl _
  · split
    · split
      · exact dropWhileLen _ _
      · exact setLength
    · exact setLength

private def trimRightMatchingBytes (input : List (BitVec 8)) (member : BitVec 8 → Bool) :
    List (BitVec 8) :=
  (input.reverse.dropWhile member).reverse

private def trimRightBySet (input cutset : List (BitVec 8)) : List (BitVec 8) :=
  let (table, ascii) := makeASCIISet cutset
  if ascii then trimRightMatchingBytes input table
  else
    let member := fun rune => containsRune (cutset.map BitVec.toNat) (Int.ofNat rune)
    let raw := input.map BitVec.toNat
    (scanTrimRightRunes member raw raw.length).map (BitVec.ofNat 8)

/-- Go 1.27.0 src/strings/strings.go:1042–1085. Remove trailing
    cutset members using byte, ASCII-set or backward UTF-8 scans. -/
def trimRight (input cutset : List (BitVec 8)) : List (BitVec 8) :=
  if input.isEmpty || cutset.isEmpty then input
  else
    match cutset with
    | [byte] => if byte.toNat < 128 then trimRightMatchingBytes input (· == byte) else trimRightBySet input cutset
    | _ => trimRightBySet input cutset

theorem trimRight_emptyCutset (input : List (BitVec 8)) :
    trimRight input [] = input := by
  simp [trimRight]

theorem trimRight_singleASCII (input : List (BitVec 8)) (byte : BitVec 8)
    (ascii : byte.toNat < 128) :
    trimRight input [byte] = (input.reverse.dropWhile (· == byte)).reverse := by
  cases input <;> simp [trimRight, trimRightMatchingBytes, ascii]

theorem trimRight_length (input cutset : List (BitVec 8)) :
    (trimRight input cutset).length ≤ input.length := by
  have byteLength (member : BitVec 8 → Bool) :
      (trimRightMatchingBytes input member).length ≤ input.length := by
    simpa [trimRightMatchingBytes] using dropWhileLen member input.reverse
  have setLength : (trimRightBySet input cutset).length ≤ input.length := by
    unfold trimRightBySet
    dsimp only
    split
    · exact byteLength _
    · simpa using scanTrimRightRunesLen _ (input.map BitVec.toNat) _
  unfold trimRight
  split
  · exact Nat.le_refl _
  · split
    · split
      · exact byteLength _
      · exact setLength
    · exact setLength

/-- Go 1.27.0 src/strings/strings.go:1027–1036. Decode the leading
    rune and drop its byte width while it occurs in the cutset. -/
def trimLeftUnicode (input cutset : List (BitVec 8)) : List (BitVec 8) :=
  let member := fun rune => containsRune (cutset.map BitVec.toNat) (Int.ofNat rune)
  let raw := input.map BitVec.toNat
  (scanTrimLeftRunes member raw raw.length).map (BitVec.ofNat 8)

theorem trimLeftUnicode_empty (cutset : List (BitVec 8)) :
    trimLeftUnicode [] cutset = [] := rfl

theorem trimLeftUnicode_length (input cutset : List (BitVec 8)) :
    (trimLeftUnicode input cutset).length ≤ input.length := by
  unfold trimLeftUnicode
  dsimp only
  simpa using scanTrimLeftRunesLen _ (input.map BitVec.toNat) _

theorem trimLeftUnicode_stop (head : BitVec 8) (tail cutset : List (BitVec 8))
    (missing : containsRune (cutset.map BitVec.toNat)
      (Int.ofNat (GoLean.Unicode.Utf8.decodeRuneInString ((head :: tail).map BitVec.toNat)).1) = false) :
    trimLeftUnicode (head :: tail) cutset = head :: tail := by
  simp [trimLeftUnicode, scanTrimLeftRunes] at missing ⊢
  have roundtrip : (BitVec.ofNat 8 ∘ BitVec.toNat) = id := by
    funext byte
    simp
  simp [missing, List.map_map, roundtrip]


/-- Go 1.27.0 src/strings/strings.go:1072–1084. Remove trailing runes
    that occur in the cutset by retaining a prefix of the original bytes. -/
def trimRightUnicode (input cutset : List (BitVec 8)) : List (BitVec 8) :=
  let member := fun rune => containsRune (cutset.map BitVec.toNat) (Int.ofNat rune)
  let raw := input.map BitVec.toNat
  (scanTrimRightRunes member raw raw.length).map (BitVec.ofNat 8)

theorem trimRightUnicode_empty (cutset : List (BitVec 8)) :
    trimRightUnicode [] cutset = [] := rfl

theorem trimRightUnicode_length (input cutset : List (BitVec 8)) :
    (trimRightUnicode input cutset).length ≤ input.length := by
  unfold trimRightUnicode
  dsimp only
  simpa using scanTrimRightRunesLen _ (input.map BitVec.toNat) _

end GoLean.Strings
