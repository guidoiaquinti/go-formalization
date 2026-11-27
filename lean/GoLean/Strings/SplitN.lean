import GoLean.Strings.Substring

namespace GoLean.Strings

private def explodeParts (input : List Nat) : Nat → List (List Nat)
  | 0 => []
  | 1 => [input]
  | fuel + 2 =>
    match input with
    | [] => []
    | lead :: rest =>
      let width := runeWidth lead rest
      input.take width :: explodeParts (input.drop width) (fuel + 1)

/-- Go 1.27.0 strings/strings.go:23–38. The last slot retains the entire suffix. -/
def explode (input : List Nat) (limit : Int) : List (List Nat) :=
  let runes := countRunes input
  explodeParts input (if limit < 0 then runes else min runes limit.toNat)

private def separatorParts (input separator : List Nat) (save : Nat) : Nat → List (List Nat)
  | 0 => [input]
  | fuel + 1 =>
    match index input separator with
    | none => [input]
    | some i =>
      input.take (i + save) :: separatorParts (input.drop (i + separator.length)) separator save fuel

/-- Go 1.27.0 strings/strings.go:274–300. None represents a nil result slice. -/
def genSplit (input separator : List Nat) (save : Nat) (limit : Int) : Option (List (List Nat)) :=
  if limit = 0 then none
  else if separator.isEmpty then some (explode input limit)
  else
    let slots := if limit < 0 then count input separator + 1 else limit.toNat
    some (separatorParts input separator save (min slots (input.length + 1) - 1))

/-- Go 1.27.0 strings/strings.go:314 passes zero retained separator bytes. -/
def splitN (input separator : List Nat) (limit : Int) : Option (List (List Nat)) :=
  genSplit input separator 0 limit

theorem splitN_zero (input separator : List Nat) : splitN input separator 0 = none := by
  simp [splitN, genSplit]

theorem splitN_one (input separator : List Nat) :
    splitN input separator 1 =
      some (if separator.isEmpty && input.isEmpty then [] else [input]) := by
  cases separator with
  | nil =>
    cases input with
    | nil => simp [splitN, genSplit, explode, explodeParts, countRunes]
    | cons lead rest =>
      simp [splitN, genSplit, explode, explodeParts, countRunes,
        Nat.min_eq_right (show 1 ≤ 1 + countRunes (rest.drop (runeWidth lead rest - 1)) by omega)]
  | cons lead rest =>
    simp [splitN, genSplit, separatorParts,
      Nat.min_eq_left (show 1 ≤ input.length + 1 by omega)]

end GoLean.Strings
