import Init.Data.BitVec.Lemmas

namespace GoLean.Strings

/-- Go 1.27.0 src/strings/strings.go:1062–1070. Scan the trailing bytes
    using the full 256-entry table. None records the nil-table panic on
    a nonempty input. -/
def trimRightASCII (input : List (BitVec 8)) (table : Option (BitVec 8 → Bool)) :
    Option (List (BitVec 8)) :=
  match input, table with
  | [], _ => some []
  | _ :: _, none => none
  | _ :: _, some member => some (input.reverse.dropWhile member).reverse

theorem trimRightASCII_empty (table : Option (BitVec 8 → Bool)) :
    trimRightASCII [] table = some [] := rfl

theorem trimRightASCII_nil (head : BitVec 8) (tail : List (BitVec 8)) :
    trimRightASCII (head :: tail) none = none := rfl

theorem trimRightASCII_retain (leading : List (BitVec 8)) (last : BitVec 8)
    (member : BitVec 8 → Bool) (missing : member last = false) :
    trimRightASCII (leading ++ [last]) (some member) = some (leading ++ [last]) := by
  cases leading <;> simp [trimRightASCII, List.reverse_append, missing]

theorem trimRightASCII_skip (leading : List (BitVec 8)) (last : BitVec 8)
    (member : BitVec 8 → Bool) (selected : member last = true) :
    trimRightASCII (leading ++ [last]) (some member) = trimRightASCII leading (some member) := by
  cases leading <;> simp [trimRightASCII, List.reverse_append, selected]

end GoLean.Strings
