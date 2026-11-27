import Init.Data.Array.BinSearch
import GoLean.Strings.Map
import GoLean.Strings.SimpleTitleData

namespace GoLean.Strings

private def scalarCount : Nat := 0x110000

private def titleTable : Array (Nat × Nat) :=
  ((simpleTitleData.splitOn " ").filterMap fun entry =>
    entry.toNat?.map fun packed => (packed / scalarCount, packed % scalarCount)).toArray

/-- Go 1.27.0 src/unicode/letter.go:289–298 and Unicode 17.0.0.
    ASCII titlecase uses uppercase arithmetic; other mappings use the pinned table. -/
def titleRune (rune : Nat) : Nat :=
  if rune < 128 then
    if 97 ≤ rune && rune ≤ 122 then rune - 32 else rune
  else
    match titleTable.binSearch (rune, 0) (fun a b => a.1 < b.1) with
    | some entry => entry.2
    | none => rune

/-- Go 1.27.0 src/strings/strings.go:747–749. Apply simple titlecase to
    each decoded rune through Map, including individual malformed bytes. -/
def toTitle (input : List (BitVec 8)) : List Nat :=
  (mapRunes (fun rune => Int.ofNat (titleRune rune)) input).2

theorem toTitle_map (input : List (BitVec 8)) :
    toTitle input = (mapRunes (fun rune => Int.ofNat (titleRune rune)) input).2 := rfl

theorem toTitle_empty : toTitle [] = [] := rfl

theorem toTitle_asciiExample : toTitle [97, 98, 32, 99, 100, 0] = [65, 66, 32, 67, 68, 0] := rfl

end GoLean.Strings
