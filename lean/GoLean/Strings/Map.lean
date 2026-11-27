import GoLean.Strings.BuilderWriteRune
import GoLean.Unicode.Utf8.DecodeRuneInString

namespace GoLean.Strings

private def scanMap (mapping : Nat → Int) (unchanged : Bool) (input : List Nat) :
    Nat → List Nat × List Nat
  | 0 => ([], [])
  | fuel + 1 =>
    match input with
    | [] => ([], [])
    | lead :: rest =>
      let (rune, width) := GoLean.Unicode.Utf8.decodeNonempty lead rest
      let mapped := mapping rune
      let keep := unchanged && mapped == (rune : Int) && (rune != 0xfffd || width != 1)
      let output := if keep then input.take width
        else if mapped < 0 then []
        else (builderRuneBytes mapped).map BitVec.toNat
      let (calls, suffix) := scanMap mapping keep (input.drop width) fuel
      (rune :: calls, output ++ suffix)

/-- Go 1.27.0 strings/strings.go:513–565. Copy the unchanged prefix,
    then encode the mapped suffix. Each decoded rune invokes the mapping
    once. An invalid byte triggers the slow path even when its RuneError
    maps to itself. -/
def mapRunes (mapping : Nat → Int) (input : List (BitVec 8)) : List Nat × List Nat :=
  scanMap mapping true (input.map BitVec.toNat) input.length

private theorem scanMap_drop (unchanged : Bool) (input : List Nat) (fuel : Nat) :
    (scanMap (fun _ => -1) unchanged input fuel).2 = [] := by
  induction fuel generalizing unchanged input with
  | zero => rfl
  | succ fuel ih =>
    cases input with
    | nil => rfl
    | cons lead rest =>
      have different : ¬ (-1 : Int) = (GoLean.Unicode.Utf8.decodeNonempty lead rest).1 := by omega
      simp [scanMap, different, ih]

theorem mapRunes_drop (input : List (BitVec 8)) :
    (mapRunes (fun _ => -1) input).2 = [] :=
  scanMap_drop _ _ _

private theorem scanMap_trace (mapping other : Nat → Int) (unchanged otherState : Bool)
    (input : List Nat) (fuel : Nat) :
    (scanMap mapping unchanged input fuel).1 = (scanMap other otherState input fuel).1 := by
  induction fuel generalizing unchanged otherState input with
  | zero => rfl
  | succ fuel ih =>
    cases input with
    | nil => rfl
    | cons lead rest =>
      simp [scanMap]
      exact ih _ _ _

/-- Mapping changes output bytes without changing the decoded-rune trace. -/
theorem mapRunes_trace (mapping other : Nat → Int) (input : List (BitVec 8)) :
    (mapRunes mapping input).1 = (mapRunes other input).1 :=
  scanMap_trace _ _ _ _ _ _

end GoLean.Strings
