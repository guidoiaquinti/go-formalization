import GoLean.Strings.SplitN

namespace GoLean.Strings

private def checkedParts (input separator : List Nat) (save : Int) :
    Nat → Except Unit (List (List Nat))
  | 0 => .ok [input]
  | fuel + 1 =>
    match index input separator with
    | none => .ok [input]
    | some i =>
      let endAt := (i : Int) + save
      if endAt < 0 || (input.length : Int) < endAt then
        .error ()
      else
        match checkedParts (input.drop (i + separator.length)) separator save fuel with
        | .error () => .error ()
        | .ok rest => .ok (input.take endAt.toNat :: rest)

/-- Go 1.27.0 strings/strings.go:274–299. The direct helper accepts a signed
    retained width. Its slice endpoint only needs to be in bounds when a
    match is processed. An error represents a slice-bounds panic. -/
def genSplitDirect (input separator : List Nat) (save limit : Int) :
    Except Unit (Option (List (List Nat))) :=
  if limit = 0 then .ok none
  else if separator.isEmpty then .ok (some (explode input limit))
  else
    let slots := if limit < 0 then count input separator + 1 else limit.toNat
    match checkedParts input separator save (min slots (input.length + 1) - 1) with
    | .error () => .error ()
    | .ok parts => .ok (some parts)

theorem genSplitDirect_zero (input separator : List Nat) (save : Int) :
    genSplitDirect input separator save 0 = .ok none := by
  simp [genSplitDirect]

theorem genSplitDirect_emptySeparator (input : List Nat) (save limit : Int) :
    genSplitDirect input [] save limit = .ok (genSplit input [] 0 limit) := by
  by_cases zero : limit = 0 <;> simp [genSplitDirect, genSplit, zero]

theorem genSplitDirect_panicFirst (input separator : List Nat) (save : Int) (i fuel : Nat)
    (nonempty : separator ≠ []) (found : index input separator = some i)
    (slots : min (fuel + 2) (input.length + 1) = fuel + 2)
    (outside : (i : Int) + save < 0 ∨ (input.length : Int) < (i : Int) + save) :
    genSplitDirect input separator save (((fuel + 2 : Nat) : Int)) = .error () := by
  have invalid : ((decide ((i : Int) + save < 0)) ||
      decide ((input.length : Int) < (i : Int) + save)) = true := by
    simpa using outside
  have rest : fuel + 2 - 1 = fuel + 1 := by omega
  have nz : (((fuel + 2 : Nat) : Int)) ≠ 0 := by omega
  have nonnegative : ¬ (((fuel + 2 : Nat) : Int)) < 0 := by omega
  have toNat : (((fuel + 2 : Nat) : Int)).toNat = fuel + 2 := rfl
  simp only [genSplitDirect, nz, nonnegative, ite_false, List.isEmpty_iff, nonempty,
    toNat, slots, rest, checkedParts, found, invalid, ite_true]

end GoLean.Strings
