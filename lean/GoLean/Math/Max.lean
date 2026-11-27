namespace GoLean.Math

private def infinityWord : Nat := 0x7ff0000000000000
private def signWord : Nat := 0x8000000000000000

private def maxRank (word : Nat) : Int :=
  if word < signWord then (word : Int) else (signWord : Int) - word

/-- Model the ordered Go branches on binary64 encodings. The public comparison
    observes NaN classification because architecture code can choose its payload.
    https://github.com/golang/go/blob/go1.27.0/src/math/dim.go#L29-L64
    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L7-L16 -/
def maxFloat64 (x y : Nat) : Nat :=
  let a := maxRank x
  let b := maxRank y
  if x = infinityWord ∨ y = infinityWord then infinityWord
  else if infinityWord < a.natAbs ∨ infinityWord < b.natAbs then 0x7ff8000000000001
  else if a = 0 ∧ a = b then
    if signWord ≤ x then y else x
  else if b < a then x else y

private def orderKey (word : Nat) : Nat :=
  if word < signWord then signWord + word else 2 ^ 64 - 1 - word

private theorem rankAbs (word : Nat) (bound : word < 2 ^ 64) :
    (maxRank word).natAbs = word % signWord := by
  unfold maxRank signWord
  split <;> omega

/-- Positive infinity takes priority over NaN. Otherwise the ordered result
    selects the greater unsigned order key, placing positive zero after negative zero. -/
theorem maxFloat64_spec (x y : Nat) (hx : x < 2 ^ 64) (hy : y < 2 ^ 64) :
    maxFloat64 x y =
      if x = infinityWord ∨ y = infinityWord then infinityWord
      else if infinityWord < x % signWord ∨ infinityWord < y % signWord then 0x7ff8000000000001
      else if orderKey y < orderKey x then x else y := by
  simp only [maxFloat64, rankAbs x hx, rankAbs y hy]
  split
  · rfl
  split
  · rfl
  by_cases sx : x < signWord <;> by_cases sy : y < signWord
  all_goals simp only [maxRank, orderKey, sx, sy, ite_true, ite_false]
  all_goals unfold signWord at *
  all_goals split <;> try split <;> try split
  all_goals omega

end GoLean.Math
