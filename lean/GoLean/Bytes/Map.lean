import GoLean.Unicode.Utf8.DecodeRune
import GoLean.Unicode.Utf8.AppendRune

namespace GoLean.Bytes

private def scanMap (mapping : Nat → Int) (input : List Nat) :
    Nat → List Nat × List Nat
  | 0 => ([], [])
  | fuel + 1 =>
    match input with
    | [] => ([], [])
    | _ :: _ =>
      let (rune, width) := GoLean.Unicode.Utf8.decodeRune input
      let mapped := mapping rune
      let output := if mapped < 0 then [] else GoLean.Unicode.Utf8.appendRune [] mapped
      let (calls, suffix) := scanMap mapping (input.drop width) fuel
      (rune :: calls, output ++ suffix)

/-- Go 1.27.0 src/bytes/bytes.go:594–611. The initial allocation is
    nonnil even when no bytes are appended. -/
def mapRunes (mapping : Nat → Int) (input : List (BitVec 8)) :
    List Nat × Option (List Nat) :=
  let (calls, bytes) := scanMap mapping (input.map BitVec.toNat) input.length
  (calls, some bytes)

private theorem scanMap_drop (input : List Nat) (fuel : Nat) :
    (scanMap (fun _ => -1) input fuel).2 = [] := by
  induction fuel generalizing input with
  | zero => rfl
  | succ fuel ih =>
    cases input with
    | nil => rfl
    | cons lead tail => simp [scanMap, ih]

private theorem scanMap_trace (f g : Nat → Int) (input : List Nat) (fuel : Nat) :
    (scanMap f input fuel).1 = (scanMap g input fuel).1 := by
  induction fuel generalizing input with
  | zero => rfl
  | succ fuel ih =>
    cases input with
    | nil => rfl
    | cons lead tail => simp [scanMap, ih]

theorem mapRunes_drop (input : List (BitVec 8)) :
    (mapRunes (fun _ => -1) input).2 = some [] := by
  simp [mapRunes, scanMap_drop]

theorem mapRunes_trace (f g : Nat → Int) (input : List (BitVec 8)) :
    (mapRunes f input).1 = (mapRunes g input).1 := by
  exact scanMap_trace f g _ _

theorem mapRunes_empty (f : Nat → Int) : mapRunes f [] = ([], some []) := rfl

end GoLean.Bytes
