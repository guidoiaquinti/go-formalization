import Init.Data.Array.BinSearch
import GoLean.Strings.Map
import GoLean.Strings.SimpleLowerData

namespace GoLean.Strings

private def scalarCount : Nat := 0x110000

private def lowerTable : Array (Nat × Nat) :=
  ((simpleLowerData.splitOn " ").filterMap fun entry =>
    entry.toNat?.map fun packed => (packed / scalarCount, packed % scalarCount)).toArray

/-- Pinned Unicode 17.0.0 simple lowercase mapping. Identity entries
    are implicit. -/
def lowerRune (rune : Nat) : Nat :=
  match lowerTable.binSearch (rune, 0) (fun a b => a.1 < b.1) with
  | some entry => entry.2
  | none => rune

private def lowerASCII (byte : Nat) : Nat :=
  if 65 ≤ byte && byte ≤ 90 then byte + 32 else byte

/-- Go 1.27.0 src/strings/strings.go:708–745. Lowercase ASCII bytes
    directly; otherwise use the existing model of Map and its treatment
    of malformed bytes. Buffer allocation is represented by output bytes. -/
def toLower (input : List (BitVec 8)) : List Nat :=
  if input.all (fun byte => byte.toNat < 128) then input.map (fun byte => lowerASCII byte.toNat)
  else (mapRunes (fun rune => Int.ofNat (lowerRune rune)) input).2

theorem toLower_ascii (input : List (BitVec 8))
    (ascii : input.all (fun byte => byte.toNat < 128) = true) :
    toLower input = input.map (fun byte => if 65 ≤ byte.toNat ∧ byte.toNat ≤ 90
      then byte.toNat + 32 else byte.toNat) := by
  simp [toLower, ascii, lowerASCII, Bool.and_eq_true]

theorem toLower_empty : toLower [] = [] := rfl

theorem toLower_asciiExample : toLower [65, 90, 97, 122, 0] = [97, 122, 97, 122, 0] := rfl

end GoLean.Strings
