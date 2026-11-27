import Init.Data.BitVec.Basic
import GoLean.Strings.Count

namespace GoLean.Strings

structure ByteStringReplacement where
  keys : List (BitVec 8)
  lookup : BitVec 8 → Option (List (BitVec 8))

private def replacementPlan (mapping : ByteStringReplacement) (input : List (BitVec 8)) :
    Int × Bool :=
  let initial := (Int.ofNat input.length, false)
  if mapping.keys.length * 8 ≤ input.length then
    let bytes := input.map BitVec.toNat
    mapping.keys.foldl (fun (size, changed) byte =>
      let occurrences := count bytes [byte.toNat]
      if occurrences = 0 then (size, changed)
      else
        let delta := Int.ofNat ((mapping.lookup byte).getD []).length - 1
        (size + Int.ofNat occurrences * delta, true)) initial
  else
    input.foldl (fun (size, changed) byte =>
      match mapping.lookup byte with
      | none => (size, changed)
      | some replacement => (size + Int.ofNat replacement.length - 1, true)) initial

/-- Go 1.27.0 src/strings/replace.go:490–548. The size pass either counts
    the listed keys or scans the input. Keys are unique and enumerate the
    present lookup entries. None at the receiver denotes a nil pointer. -/
def byteStringReplacerReplace (mapping : Option ByteStringReplacement)
    (input : List (BitVec 8)) : Option (List (BitVec 8)) := do
  let mapping ← mapping
  let (size, changed) := replacementPlan mapping input
  if !changed then return input
  let output := input.flatMap fun byte => (mapping.lookup byte).getD [byte]
  if size = Int.ofNat output.length then return output
  else none

theorem byteStringReplacerReplace_nil (input : List (BitVec 8)) :
    byteStringReplacerReplace none input = none := rfl

theorem byteStringReplacerReplace_empty (mapping : ByteStringReplacement) :
    byteStringReplacerReplace (some mapping) [] = some [] := by
  cases mapping with
  | mk keys lookup => cases keys <;> simp [byteStringReplacerReplace, replacementPlan]

theorem byteStringReplacerReplace_noKeys (lookup : BitVec 8 → Option (List (BitVec 8)))
    (input : List (BitVec 8)) :
    byteStringReplacerReplace (some ⟨[], lookup⟩) input = some input := by
  simp [byteStringReplacerReplace, replacementPlan]

theorem byteStringReplacerReplace_noncascading :
    let byteMap (byte : BitVec 8) :=
      if byte = 97 then some [98, 0xff] else if byte = 98 then some [] else none
    byteStringReplacerReplace (some ⟨[97, 98], byteMap⟩) [97, 98, 0xff] =
      some [98, 0xff, 0xff] := by rfl

end GoLean.Strings
