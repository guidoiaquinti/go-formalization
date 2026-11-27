import GoLean.Strings.SplitN

namespace GoLean.Strings

/-- The shared source-shaped model used by SplitN's empty-separator branch.
    The result list always represents a nonnil Go slice.
    Go 1.27.0 strings/strings.go:23–38. -/
abbrev explodeDirect := explode

theorem explode_zero (input : List Nat) : explodeDirect input 0 = [] := by
  simp [explodeDirect, explode]
  rfl

theorem explode_negative (input : List Nat) (limit : Int) (negative : limit < 0) :
    explodeDirect input limit = explodeDirect input (-1) := by
  simp [explodeDirect, explode, negative]

theorem explode_one (input : List Nat) :
    explodeDirect input 1 = (if input.isEmpty then [] else [input]) := by
  cases input with
  | nil => simp [explodeDirect, explode, countRunes]; rfl
  | cons lead rest =>
    simp [explodeDirect, explode, countRunes, Nat.min_eq_right
      (show 1 ≤ 1 + countRunes (rest.drop (runeWidth lead rest - 1)) by omega)]
    rfl

end GoLean.Strings
