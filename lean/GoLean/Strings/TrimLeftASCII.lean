import Init.Data.BitVec.Lemmas

namespace GoLean.Strings

/-- Go 1.27.0 src/strings/strings.go:1017–1025. Index the entire
    256-entry table with each leading byte. None records the nil-table
    panic when a nonempty input requires its first lookup. -/
def trimLeftASCII (input : List (BitVec 8)) (table : Option (BitVec 8 → Bool)) :
    Option (List (BitVec 8)) :=
  match input, table with
  | [], _ => some []
  | _ :: _, none => none
  | _ :: _, some member => some (input.dropWhile member)

theorem trimLeftASCII_empty (table : Option (BitVec 8 → Bool)) :
    trimLeftASCII [] table = some [] := rfl

theorem trimLeftASCII_nil (head : BitVec 8) (tail : List (BitVec 8)) :
    trimLeftASCII (head :: tail) none = none := rfl

theorem trimLeftASCII_retain (head : BitVec 8) (tail : List (BitVec 8))
    (member : BitVec 8 → Bool) (missing : member head = false) :
    trimLeftASCII (head :: tail) (some member) = some (head :: tail) := by
  simp [trimLeftASCII, List.dropWhile, missing]

theorem trimLeftASCII_skip (head : BitVec 8) (tail : List (BitVec 8))
    (member : BitVec 8 → Bool) (selected : member head = true) :
    trimLeftASCII (head :: tail) (some member) = trimLeftASCII tail (some member) := by
  cases tail <;> simp [trimLeftASCII, List.dropWhile, selected]

end GoLean.Strings
