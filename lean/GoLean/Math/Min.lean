namespace GoLean.Math

private def infinityWord : Nat := 0x7ff0000000000000
private def signWord : Nat := 0x8000000000000000

private def minRank (word : Nat) : Int :=
  if word < signWord then (word : Int) else (signWord : Int) - word

/-- Model the ordered Go branches on binary64 encodings. The public comparison
    observes NaN classification because architecture code can choose its payload.
    https://github.com/golang/go/blob/go1.27.0/src/math/dim.go#L66-L100
    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L7-L16 -/
def minFloat64 (x y : Nat) : Nat :=
  let a := minRank x
  let b := minRank y
  if x = signWord + infinityWord ∨ y = signWord + infinityWord then signWord + infinityWord
  else if infinityWord < a.natAbs ∨ infinityWord < b.natAbs then 0x7ff8000000000001
  else if a = 0 ∧ a = b then
    if signWord ≤ x then x else y
  else if a < b then x else y

private def orderKey (word : Nat) : Nat :=
  if word < signWord then signWord + word else 2 ^ 64 - 1 - word

private theorem rankAbs (word : Nat) (bound : word < 2 ^ 64) :
    (minRank word).natAbs = word % signWord := by
  unfold minRank signWord
  split <;> omega

/-- Negative infinity takes priority over NaN. Otherwise the ordered result
    selects the smaller unsigned order key, placing negative zero before positive zero. -/
theorem minFloat64_spec (x y : Nat) (hx : x < 2 ^ 64) (hy : y < 2 ^ 64) :
    minFloat64 x y =
      if x = signWord + infinityWord ∨ y = signWord + infinityWord then signWord + infinityWord
      else if infinityWord < x % signWord ∨ infinityWord < y % signWord then 0x7ff8000000000001
      else if orderKey x < orderKey y then x else y := by
  simp only [minFloat64, rankAbs x hx, rankAbs y hy]
  split
  · rfl
  split
  · rfl
  by_cases sx : x < signWord <;> by_cases sy : y < signWord
  all_goals simp only [minRank, orderKey, sx, sy, ite_true, ite_false]
  all_goals unfold signWord at *
  all_goals split <;> try split <;> try split
  all_goals omega

end GoLean.Math
