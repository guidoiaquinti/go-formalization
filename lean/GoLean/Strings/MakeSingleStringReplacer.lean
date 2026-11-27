import GoLean.Strings.MakeStringFinder

namespace GoLean.Strings

structure SingleStringReplacer where
  finder : FinderTables
  value : List (BitVec 8)

/-- Go 1.27.0 src/strings/replace.go:385–387. Store the finder built for
    the original pattern and the replacement bytes verbatim. -/
def makeSingleStringReplacer (pattern value : List (BitVec 8)) : SingleStringReplacer :=
  ⟨makeStringFinder pattern, value⟩

theorem makeSingleStringReplacer_fields (pattern value : List (BitVec 8)) :
    (makeSingleStringReplacer pattern value).finder = makeStringFinder pattern ∧
    (makeSingleStringReplacer pattern value).value = value :=
  ⟨rfl, rfl⟩

end GoLean.Strings
