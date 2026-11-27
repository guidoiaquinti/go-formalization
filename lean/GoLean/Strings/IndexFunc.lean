import GoLean.Unicode.Utf8.DecodeRuneInString

namespace GoLean.Strings

private def scanIndexFunc (predicate : Nat → Bool) (truth : Bool)
    (offset : Nat) (input : List Nat) : Nat → Int × List Nat
  | 0 => (-1, [])
  | fuel + 1 =>
    match input with
    | [] => (-1, [])
    | lead :: rest =>
      let (rune, width) := GoLean.Unicode.Utf8.decodeNonempty lead rest
      if predicate rune == truth then (offset, [rune])
      else
        let (index, calls) := scanIndexFunc predicate truth (offset + width) (input.drop width) fuel
        (index, rune :: calls)

/-- Go 1.27.0 strings/strings.go:910–917. None is the nil callback
    panic. Some includes the byte index and callback trace. -/
def indexFuncWithTruth (input : List (BitVec 8)) (predicate : Option (Nat → Bool))
    (truth : Bool) : Option (Int × List Nat) :=
  match input, predicate with
  | [], _ => some (-1, [])
  | _ :: _, none => none
  | _ :: _, some f => some (scanIndexFunc f truth 0 (input.map BitVec.toNat) input.length)

/-- Go 1.27.0 strings/strings.go:897–899 searches for a true result. -/
def indexFunc (input : List (BitVec 8)) (predicate : Option (Nat → Bool)) :
    Option (Int × List Nat) :=
  indexFuncWithTruth input predicate true

theorem indexFunc_empty (predicate : Option (Nat → Bool)) :
    indexFunc [] predicate = some (-1, []) := rfl

theorem indexFunc_nil (lead : BitVec 8) (rest : List (BitVec 8)) :
    indexFunc (lead :: rest) none = none := rfl

theorem indexFunc_first (lead : BitVec 8) (rest : List (BitVec 8)) (f : Nat → Bool)
    (matched : f (GoLean.Unicode.Utf8.decodeNonempty lead.toNat (rest.map BitVec.toNat)).1 = true) :
    indexFunc (lead :: rest) (some f) =
      some (0, [(GoLean.Unicode.Utf8.decodeNonempty lead.toNat (rest.map BitVec.toNat)).1]) := by
  simp [indexFunc, indexFuncWithTruth, scanIndexFunc, matched]

private theorem scanIndexFunc_false (input : List Nat) (offset fuel : Nat) :
    (scanIndexFunc (fun _ => false) true offset input fuel).1 = -1 := by
  induction fuel generalizing input offset with
  | zero => rfl
  | succ fuel ih =>
    cases input with
    | nil => rfl
    | cons lead rest => simp [scanIndexFunc, ih]

theorem indexFunc_noMatch (input : List (BitVec 8)) :
    (indexFunc input (some (fun _ => false))).map Prod.fst = some (-1) := by
  cases input with
  | nil => rfl
  | cons lead rest =>
    change some (scanIndexFunc (fun _ => false) true 0 ((lead :: rest).map BitVec.toNat)
      (lead :: rest).length).1 = _
    rw [scanIndexFunc_false]

end GoLean.Strings
