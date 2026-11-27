import Init.Data.Array.BinSearch
import GoLean.Strings.SimpleFoldData
import GoLean.Unicode.Utf8.DecodeRuneInString

namespace GoLean.Strings

private def scalarCount : Nat := 0x110000

private def foldTable : Array (Nat × Nat) :=
  ((simpleFoldData.splitOn " ").filterMap fun entry =>
    entry.toNat?.map fun packed => (packed / scalarCount, packed % scalarCount)).toArray

/-- Unicode 17.0.0 simple-fold successor from the exact Go 1.27.0 table.
    Identity transitions are implicit. -/
def equalFoldSuccessor (rune : Nat) : Nat :=
  match foldTable.binSearch (rune, 0) (fun a b => a.1 < b.1) with
  | some entry => entry.2
  | none => rune

private def asciiFoldPair (small large : Nat) : Bool :=
  small == large || (65 ≤ small && small ≤ 90 && large == small + 32)

private def followFold (start target current : Nat) : Nat → Bool
  | 0 => false
  | fuel + 1 =>
    if current = target then true
    else if current = start || target < current then false
    else followFold start target (equalFoldSuccessor current) fuel

private def foldPair (a b : Nat) : Bool :=
  let small := min a b
  let large := max a b
  if small = large then true
  else if large < 128 then asciiFoldPair small large
  else followFold small large (equalFoldSuccessor small) 4

private def equalFoldRunes (a b : List Nat) : Nat → Bool
  | 0 => a.isEmpty && b.isEmpty
  | fuel + 1 =>
    match a, b with
    | [], [] => true
    | [], _ | _, [] => false
    | _ :: _, _ :: _ =>
      let (first, aw) := GoLean.Unicode.Utf8.decodeRuneInString a
      let (second, bw) := GoLean.Unicode.Utf8.decodeRuneInString b
      foldPair first second && equalFoldRunes (a.drop aw) (b.drop bw) fuel

/-- Go 1.27.0 src/strings/strings.go:1185–1260. Scan ASCII bytes first.
    When either byte is non-ASCII, decode one rune on each side and follow
    the sorted SimpleFold orbit. Malformed bytes consume one byte each. -/
def equalFold : List Nat → List Nat → Bool
  | [], [] => true
  | [], _ | _, [] => false
  | a :: as, b :: bs =>
    if a < 128 && b < 128 then
      asciiFoldPair (min a b) (max a b) && equalFold as bs
    else equalFoldRunes (a :: as) (b :: bs) (as.length + 1)

theorem equalFold_asciiHead (a b : Nat) (as bs : List Nat) (asciiA : a < 128) (asciiB : b < 128) :
    equalFold (a :: as) (b :: bs) = true ↔
      (a = b ∨ (65 ≤ a ∧ a ≤ 90 ∧ b = a + 32) ∨ (65 ≤ b ∧ b ≤ 90 ∧ a = b + 32)) ∧
        equalFold as bs = true := by
  simp only [equalFold, asciiA, asciiB, decide_true, Bool.and_self, ite_true,
    asciiFoldPair, Bool.and_eq_true, Bool.or_eq_true, beq_iff_eq, decide_eq_true_eq]
  by_cases ab : a ≤ b
  · rw [Nat.min_eq_left ab, Nat.max_eq_right ab]
    by_cases tails : equalFold as bs = true <;> simp [tails] <;> omega
  · have ba : b ≤ a := by omega
    rw [Nat.min_eq_right ba, Nat.max_eq_left ba]
    by_cases tails : equalFold as bs = true <;> simp [tails] <;> omega

theorem equalFold_invalidRuneError :
    equalFold [0xff] [0xef, 0xbf, 0xbd] = true := by rfl

theorem equalFold_invalidByteCount :
    equalFold [0xff, 0x80] [0xef, 0xbf, 0xbd] = false := by rfl

end GoLean.Strings
