import Init.Data.Array.BinSearch
import GoLean.Strings.Map
import GoLean.Strings.SimpleUpperData

namespace GoLean.Strings

private def scalarCount : Nat := 0x110000

private def upperTable : Array (Nat × Nat) :=
  ((simpleUpperData.splitOn " ").filterMap fun entry =>
    entry.toNat?.map fun packed => (packed / scalarCount, packed % scalarCount)).toArray

/-- Pinned Unicode 17.0.0 simple uppercase mapping. Identity entries
    are implicit. -/
def upperRune (rune : Nat) : Nat :=
  match upperTable.binSearch (rune, 0) (fun a b => a.1 < b.1) with
  | some entry => entry.2
  | none => rune

private def upperASCII (byte : Nat) : Nat :=
  if 97 ≤ byte && byte ≤ 122 then byte - 32 else byte

/-- Go 1.27.0 src/strings/strings.go:668–705. Uppercase ASCII bytes
    directly; otherwise use the existing model of Map and its treatment
    of malformed bytes. Buffer allocation is represented by output bytes. -/
def toUpper (input : List (BitVec 8)) : List Nat :=
  if input.all (fun byte => byte.toNat < 128) then input.map (fun byte => upperASCII byte.toNat)
  else (mapRunes (fun rune => Int.ofNat (upperRune rune)) input).2

theorem toUpper_ascii (input : List (BitVec 8))
    (ascii : input.all (fun byte => byte.toNat < 128) = true) :
    toUpper input = input.map (fun byte => if 97 ≤ byte.toNat ∧ byte.toNat ≤ 122
      then byte.toNat - 32 else byte.toNat) := by
  simp [toUpper, ascii, upperASCII, Bool.and_eq_true]

theorem toUpper_empty : toUpper [] = [] := rfl

theorem toUpper_asciiExample : toUpper [65, 90, 97, 122, 0] = [65, 90, 65, 90, 0] := rfl

end GoLean.Strings
