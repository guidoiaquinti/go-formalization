import GoLean.Strings.IsSeparator
import GoLean.Strings.ToTitle

namespace GoLean.Strings

private def scanTitle (previous : Nat) (unchanged : Bool) (input : List Nat) :
    Nat → List Nat
  | 0 => []
  | fuel + 1 =>
    match input with
    | [] => []
    | lead :: rest =>
      let (rune, width) := GoLean.Unicode.Utf8.decodeNonempty lead rest
      let mapped := if isSeparator (Int.ofNat previous) then titleRune rune else rune
      let keep := unchanged && mapped == rune && (rune != 0xfffd || width != 1)
      let output := if keep then input.take width
        else (builderRuneBytes (Int.ofNat mapped)).map BitVec.toNat
      output ++ scanTitle rune keep (input.drop width) fuel

/-- Go 1.27.0 src/strings/strings.go:844–864. Map calls the closure
    once per decoded rune, keeping the previous original rune as state. -/
def title (input : List (BitVec 8)) : List Nat :=
  scanTitle 32 true (input.map BitVec.toNat) input.length

theorem title_empty : title [] = [] := rfl

theorem title_asciiSeparators :
    title [97, 98, 45, 99, 100, 95, 101, 102, 32, 103] =
      [65, 98, 45, 67, 100, 95, 101, 102, 32, 71] := rfl

theorem title_asciiDigit : title [49, 97, 46, 98] = [49, 97, 46, 66] := rfl

theorem title_asciiNUL : title [97, 0, 98] = [65, 0, 66] := rfl

end GoLean.Strings
