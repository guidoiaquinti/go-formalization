import GoLean.Bytes.SplitN

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:34–51. Direct calls treat a zero bound
    as unbounded. The empty-separator scanner caps decoded fragments
    and preserves any spare capacity in the last raw suffix. -/
def explode (input : List Nat) (spare : Nat) (limit : Int) : Option (List (List Nat × Nat)) :=
  splitN input [] spare (if limit ≤ 0 then -1 else limit)

theorem explode_nonpositive (input : List Nat) (spare : Nat) (limit : Int)
    (nonpositive : limit ≤ 0) :
    explode input spare limit = explode input spare (-1) := by
  simp [explode, nonpositive]

theorem explode_one (lead : Nat) (rest : List Nat) (spare : Nat) :
    explode (lead :: rest) spare 1 =
      some [(lead :: rest, (lead :: rest).length + spare)] := by
  simpa [explode] using splitN_one lead rest [] spare

theorem explode_empty (spare : Nat) (limit : Int) :
    explode [] spare limit = some [] := by
  have h : (if limit ≤ 0 then (-1 : Int) else limit) ≠ 0 := by split <;> omega
  simp [explode, splitN, withSplitCaps, GoLean.Strings.splitN, GoLean.Strings.genSplit,
    GoLean.Strings.explode, GoLean.Strings.countRunes, h]
  rfl

end GoLean.Bytes
