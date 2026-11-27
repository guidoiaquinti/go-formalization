import GoLean.Unicode.Utf8.DecodeLastRuneInString

namespace GoLean.Strings

private def scanLastIndexFunc (predicate : Nat → Bool) (truth : Bool)
    (input : List Nat) : Nat → Int × List Nat
  | 0 => (-1, [])
  | fuel + 1 =>
    if input.isEmpty then (-1, [])
    else
      let (rune, width) := GoLean.Unicode.Utf8.decodeLastRuneInString input
      let offset := input.length - width
      if predicate rune == truth then (offset, [rune])
      else
        let (index, calls) := scanLastIndexFunc predicate truth (input.take offset) fuel
        (index, rune :: calls)

/-- Go 1.27.0 strings/strings.go:922–933. None represents the nil callback
    panic. The index is the byte offset and the trace is in call order. -/
def lastIndexFuncWithTruth (input : List (BitVec 8)) (predicate : Option (Nat → Bool))
    (truth : Bool) : Option (Int × List Nat) :=
  match input, predicate with
  | [], _ => some (-1, [])
  | _ :: _, none => none
  | _ :: _, some f => some (scanLastIndexFunc f truth (input.map BitVec.toNat) input.length)

/-- Go 1.27.0 strings/strings.go:903–905 searches backward for a true result. -/
def lastIndexFunc (input : List (BitVec 8)) (predicate : Option (Nat → Bool)) :
    Option (Int × List Nat) :=
  lastIndexFuncWithTruth input predicate true

theorem lastIndexFunc_empty (predicate : Option (Nat → Bool)) :
    lastIndexFunc [] predicate = some (-1, []) := rfl

theorem lastIndexFunc_nil (head : BitVec 8) (tail : List (BitVec 8)) :
    lastIndexFunc (head :: tail) none = none := rfl

theorem lastIndexFunc_firstCall (head : BitVec 8) (tail : List (BitVec 8)) (f : Nat → Bool)
    (matched : f (GoLean.Unicode.Utf8.decodeLastRuneInString ((head :: tail).map BitVec.toNat)).1 = true) :
    lastIndexFunc (head :: tail) (some f) =
      let decoded := GoLean.Unicode.Utf8.decodeLastRuneInString ((head :: tail).map BitVec.toNat)
      some ((Int.ofNat ((head :: tail).length - decoded.2)), [decoded.1]) := by
  simp [lastIndexFunc, lastIndexFuncWithTruth, scanLastIndexFunc] at matched ⊢
  simp [matched]

end GoLean.Strings
