import GoLean.Strings.GenericReplacerWriteString

namespace GoLean.Strings

/-- Go 1.27.0 src/strings/replace.go:330–334. The append-slice writer
    consumes every emitted byte span without an error. The returned string
    contains those spans in order. -/
def genericReplacerReplace (pairs : List (List (BitVec 8) × List (BitVec 8)))
    (input : List (BitVec 8)) : List (BitVec 8) :=
  (genericReplacerWriteSegments pairs input).flatten

theorem genericReplacerReplace_empty : genericReplacerReplace [] [] = [] := rfl

theorem genericReplacerReplace_emptyKey (replacement : List (BitVec 8)) :
    genericReplacerReplace [([], replacement)] [0xff] =
      replacement ++ [0xff] ++ replacement := by
  change replacement ++ (0xff :: replacement ++ []) = replacement ++ [0xff] ++ replacement
  simp

theorem genericReplacerReplace_nonoverlap (replacement : List (BitVec 8)) :
    genericReplacerReplace [([97, 97], replacement)] [97, 97, 97] =
      replacement ++ [97] := by
  change replacement ++ [97] = replacement ++ [97]
  rfl

theorem genericReplacerReplace_priority (first later : List (BitVec 8)) :
    genericReplacerReplace [([97], first), ([97, 98], later)] [97, 98] = first ++ [98] := rfl

theorem genericReplacerReplace_rootThenKey (root key : List (BitVec 8)) :
    genericReplacerReplace [([], root), ([0xff], key)] [0xff] = root ++ key ++ root := by
  change root ++ (key ++ (root ++ [])) = root ++ key ++ root
  simp [List.append_assoc]

end GoLean.Strings
