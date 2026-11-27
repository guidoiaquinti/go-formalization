import GoLean.Strings.IndexRune

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:131–203. The result is a byte offset
    or None for the Go sentinel -1. Candidate probes and cutover
    implementations are abstracted in the shared result model. -/
def indexRune (input : List Nat) (r : Int) : Option Nat :=
  GoLean.Strings.indexRune input r

theorem indexRune_valid (input : List Nat) (r : Int)
    (valid : GoLean.Strings.validRune r = true) (ordinary : r ≠ 0xfffd) :
    GoLean.Strings.FirstMatch input (GoLean.Strings.encodeRune r.toNat) (indexRune input r) :=
  GoLean.Strings.indexRune_valid input r valid ordinary

theorem indexRune_invalid (input : List Nat) (r : Int)
    (invalid : GoLean.Strings.validRune r = false) (ordinary : r ≠ 0xfffd) :
    indexRune input r = none :=
  GoLean.Strings.indexRune_invalid input r invalid ordinary

theorem indexRune_error_cons (lead : Nat) (rest : List Nat) :
    indexRune (lead :: rest) 0xfffd =
      if GoLean.Strings.runeErrorHere lead rest then some 0
      else (indexRune (rest.drop (GoLean.Strings.runeWidth lead rest - 1)) 0xfffd).map
        (fun i => GoLean.Strings.runeWidth lead rest + i) :=
  GoLean.Strings.indexRune_error_cons lead rest

end GoLean.Bytes
