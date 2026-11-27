import GoLean.Bytes.SplitN

namespace GoLean.Bytes

private def checkedParts (input separator spare : List Nat) (save : Int) :
    Nat → Except Unit (List (List Nat × Nat))
  | 0 => .ok [(input, input.length + spare.length)]
  | fuel + 1 =>
    match GoLean.Strings.index input separator with
    | none => .ok [(input, input.length + spare.length)]
    | some i =>
      let endAt := (i : Int) + save
      if endAt < 0 || ((input.length + spare.length : Nat) : Int) < endAt then
        .error ()
      else do
        let rest ← checkedParts (input.drop (i + separator.length)) separator spare save fuel
        return ((input ++ spare).take endAt.toNat, endAt.toNat) :: rest

/-- Go 1.27.0 bytes/bytes.go:373–397. The full slice expression can reach
    the initialized backing bytes between length and capacity. -/
def genSplit (input separator spare : List Nat) (save limit : Int) :
    Except Unit (Option (List (List Nat × Nat))) :=
  if limit = 0 then .ok none
  else if separator.isEmpty then .ok (splitN input separator spare.length limit)
  else
    let slots := if limit < 0 then GoLean.Strings.count input separator + 1 else limit.toNat
    (checkedParts input separator spare save (min slots (input.length + 1) - 1)).map some

theorem genSplit_zero (input separator spare : List Nat) (save : Int) :
    genSplit input separator spare save 0 = .ok none := by
  simp [genSplit]

theorem genSplit_emptySeparator (input spare : List Nat) (save limit : Int) :
    genSplit input [] spare save limit = .ok (splitN input [] spare.length limit) := by
  by_cases zero : limit = 0 <;> simp [genSplit, zero, splitN_zero]

theorem genSplit_one (input separator spare : List Nat) (save : Int)
    (nonempty : separator ≠ []) :
    genSplit input separator spare save 1 = .ok (some [(input, input.length + spare.length)]) := by
  simp [genSplit, List.isEmpty_iff, nonempty, checkedParts, Except.map,
    Nat.min_eq_left (show 1 ≤ input.length + 1 by omega)]

end GoLean.Bytes
