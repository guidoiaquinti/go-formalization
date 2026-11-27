import Init.Data.BitVec.Basic

namespace GoLean.Strings

/-- Go 1.27.0 src/strings/replace.go:438–455. Look up each original
    byte once for the returned bytes. The lazy copy is represented by its
    result. A nil receiver is only dereferenced when the input has a byte. -/
def byteReplacerReplace (mapping : Option (BitVec 8 → BitVec 8))
    (input : List (BitVec 8)) : Option (List (BitVec 8)) :=
  match input, mapping with
  | [], _ => some []
  | _ :: _, none => none
  | _ :: _, some lookup => some (input.map lookup)

theorem byteReplacerReplace_empty (mapping : Option (BitVec 8 → BitVec 8)) :
    byteReplacerReplace mapping [] = some [] := rfl

theorem byteReplacerReplace_length (mapping : BitVec 8 → BitVec 8) (input : List (BitVec 8)) :
    (byteReplacerReplace (some mapping) input).map List.length = some input.length := by
  cases input <;> simp [byteReplacerReplace]

theorem byteReplacerReplace_head (mapping : BitVec 8 → BitVec 8) (head : BitVec 8)
    (tail : List (BitVec 8)) :
    (byteReplacerReplace (some mapping) (head :: tail)).map List.head? = some (some (mapping head)) :=
  rfl

theorem byteReplacerReplace_identity (input : List (BitVec 8)) :
    byteReplacerReplace (some id) input = some input := by
  cases input <;> simp [byteReplacerReplace]

theorem byteReplacerReplace_nil (head : BitVec 8) (tail : List (BitVec 8)) :
    byteReplacerReplace none (head :: tail) = none := rfl

end GoLean.Strings
