import GoLean.Strings.IndexRune

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:87–89, 131–200.
    The result records whether IndexRune finds the requested rune.
    The shared result model covers ASCII, RuneError, invalid runes, and
    encoded-byte matches without modeling architecture-specific probes. -/
def containsRune (input : List Nat) (r : Int) : Bool :=
  (GoLean.Strings.indexRune input r).isSome

theorem containsRune_spec (input : List Nat) (r : Int) :
    containsRune input r = true ↔ ∃ i, GoLean.Strings.indexRune input r = some i :=
  GoLean.Strings.containsRune_spec input r

theorem containsRune_invalid (input : List Nat) (r : Int)
    (invalid : GoLean.Strings.validRune r = false) (ordinary : r ≠ 0xfffd) :
    containsRune input r = false := by
  simp [containsRune, GoLean.Strings.indexRune_invalid input r invalid ordinary]

theorem containsRune_error (lead : Nat) (rest : List Nat)
    (errorHere : GoLean.Strings.runeErrorHere lead rest = true) :
    containsRune (lead :: rest) 0xfffd = true := by
  simp [containsRune, GoLean.Strings.indexRune_error_cons, errorHere]

end GoLean.Bytes
