import GoLean.Strings.CutPrefix
import GoLean.Strings.LongestCommonSuffix

namespace GoLean.Strings

structure FinderTables where
  pattern : List (BitVec 8)
  badCharSkip : List Nat
  goodSuffixSkip : List Nat

private def badCharTable (pattern : List (BitVec 8)) (last : Nat) : List Nat :=
  (List.range last).foldl (fun table i =>
    table.set (pattern[i]!).toNat (last - i)) (List.replicate 256 pattern.length)

private def prefixPass (pattern : List (BitVec 8)) (last : Nat) : List Nat :=
  ((List.range pattern.length).reverse.foldl (fun (nextPrefix, table) i =>
    let nextPrefix := if hasPrefix pattern (pattern.drop (i + 1)) then i + 1 else nextPrefix
    (nextPrefix, (nextPrefix + last - i) :: table)) (last, [])).2

private def suffixPass (pattern : List (BitVec 8)) (last : Nat) (table : List Nat) : List Nat :=
  (List.range last).foldl (fun table i =>
    let width := longestCommonSuffix pattern ((pattern.take (i + 1)).drop 1)
    if pattern[i - width]! != pattern[last - width]! then
      table.set (last - width) (width + last - i)
    else table) table

/-- Go 1.27.0 src/strings/search.go:48–89. Initialize bad-character
    distances, build prefix fallbacks, then overwrite with repeated-suffix
    shifts when the preceding byte can differ. Empty index ranges are skipped. -/
def makeStringFinder (pattern : List (BitVec 8)) : FinderTables :=
  let last := pattern.length - 1
  ⟨pattern, badCharTable pattern last, suffixPass pattern last (prefixPass pattern last)⟩

private theorem foldLength (step : List Nat → Nat → List Nat)
    (preserves : ∀ table i, (step table i).length = table.length)
    (indices table : List Nat) : (indices.foldl step table).length = table.length := by
  induction indices generalizing table with
  | nil => rfl
  | cons i indices ih => simp [List.foldl, ih, preserves]

private theorem prefixPassLength (pattern : List (BitVec 8)) (last : Nat) :
    (prefixPass pattern last).length = pattern.length := by
  unfold prefixPass
  suffices ∀ (indices : List Nat) (initialPrefix : Nat) (table : List Nat),
      ((indices.foldl (fun (nextPrefix, table) i =>
        let nextPrefix := if hasPrefix pattern (pattern.drop (i + 1)) then i + 1 else nextPrefix
        (nextPrefix, (nextPrefix + last - i) :: table)) (initialPrefix, table)).2).length =
      table.length + indices.length by
    simpa using this ((List.range pattern.length).reverse) last []
  intro indices
  induction indices with
  | nil => simp
  | cons i indices ih =>
    intro initialPrefix table
    simp [List.foldl, ih]
    omega

theorem makeStringFinder_shape (pattern : List (BitVec 8)) :
    (makeStringFinder pattern).pattern = pattern ∧
    (makeStringFinder pattern).badCharSkip.length = 256 ∧
    (makeStringFinder pattern).goodSuffixSkip.length = pattern.length := by
  refine ⟨rfl, ?_, ?_⟩
  · unfold makeStringFinder badCharTable
    rw [foldLength]
    · simp only [List.length_replicate]
    · intros
      simp
  · unfold makeStringFinder suffixPass
    rw [foldLength]
    · exact prefixPassLength _ _
    · intro table i
      dsimp only
      split <;> simp

theorem makeStringFinder_empty :
    (makeStringFinder []).badCharSkip = List.replicate 256 0 ∧
    (makeStringFinder []).goodSuffixSkip = [] := by
  exact ⟨rfl, rfl⟩

end GoLean.Strings
