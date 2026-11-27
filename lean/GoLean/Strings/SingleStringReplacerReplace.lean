import GoLean.Strings.MakeSingleStringReplacer
import GoLean.Strings.StringFinderNext

namespace GoLean.Strings

private def replaceSingle (replacer : SingleStringReplacer) (input : List (BitVec 8)) :
    Nat → List (BitVec 8)
  | 0 => input
  | fuel + 1 =>
    let matchIndex := replacer.finder.next input
    if matchIndex < 0 then input
    else
      let start := matchIndex.toNat
      input.take start ++ replacer.value ++
        replaceSingle replacer (input.drop (start + replacer.finder.pattern.length)) fuel

/-- Go 1.27.0 src/strings/replace.go:389–409, for nonempty patterns.
    Copy each unmatched segment and replacement, then resume after the
    matched bytes. The builder and unchanged-input fast path are represented
    by their returned bytes. -/
def singleStringReplacerReplace (pattern value input : List (BitVec 8)) : List (BitVec 8) :=
  replaceSingle (makeSingleStringReplacer pattern value) input (input.length + 1)

theorem singleStringReplacerReplace_noMatch (pattern value input : List (BitVec 8))
    (missing : stringFinderNext pattern input = -1) :
    singleStringReplacerReplace pattern value input = input := by
  simp [singleStringReplacerReplace, replaceSingle, makeSingleStringReplacer,
    FinderTables.next_constructed, missing]

theorem singleStringReplacerReplace_tooShort (pattern value input : List (BitVec 8))
    (short : input.length < pattern.length) :
    singleStringReplacerReplace pattern value input = input :=
  singleStringReplacerReplace_noMatch _ _ _ (stringFinderNext_tooShort _ _ short)

theorem singleStringReplacerReplace_overlap (value : List (BitVec 8)) :
    singleStringReplacerReplace [0xff, 0xff] value [0xff, 0xff, 0xff] = value ++ [0xff] := by
  rfl

end GoLean.Strings
