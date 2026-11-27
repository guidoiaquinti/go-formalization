import Init.Data.Nat.Bitwise.Lemmas

namespace GoLean.Math

/-- The Go implementation handles special values before normalizing finite
    nonzero inputs. Multiplying a subnormal by 2^52 gives a biased exponent
    of one plus the position of its highest set fraction bit.
    https://github.com/golang/go/blob/go1.27.0/src/math/logb.go#L27-L52
    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L61-L67 -/
def ilogb (word : Nat) : Int :=
  let e := (word >>> 52) &&& 2047
  let f := word &&& (2 ^ 52 - 1)
  if e = 0 ∧ f = 0 then -2147483648
  else if e = 2047 then 2147483647
  else
    let (normalized, adjustment) := if e = 0 then (f.log2 + 1, (-52 : Int)) else (e, (0 : Int))
    (normalized : Int) - 1023 + adjustment

/-- The magnitude word gives the binary exponent. A subnormal's exponent
    is the index of its highest set fraction bit less 1074. -/
theorem ilogb_spec (word : Nat) :
    ilogb word =
      let magnitude := word % 2 ^ 63
      let exponent : Nat := magnitude / 2 ^ 52
      if magnitude = 0 then -2147483648
      else if exponent = 2047 then 2147483647
      else if exponent = 0 then (magnitude.log2 : Int) - 1074
      else (exponent : Int) - 1023 := by
  unfold ilogb
  dsimp only
  rw [Nat.shiftRight_eq_div_pow, show (2047 : Nat) = 2 ^ 11 - 1 by decide,
    Nat.and_two_pow_sub_one_eq_mod, Nat.and_two_pow_sub_one_eq_mod]
  have hExp : word / 2 ^ 52 % 2 ^ 11 = word % 2 ^ 63 / 2 ^ 52 := by omega
  rw [hExp]
  have hZero : (word % 2 ^ 63 / 2 ^ 52 = 0 ∧ word % 2 ^ 52 = 0) ↔
      word % 2 ^ 63 = 0 := by omega
  by_cases zero : word % 2 ^ 63 = 0
  · rw [ite_eq_left (hZero.mpr zero), ite_eq_left zero]
  · rw [ite_eq_right (mt hZero.mp zero), ite_eq_right zero]
    change (if word % 2 ^ 63 / 2 ^ 52 = 2047 then 2147483647 else _) =
      (if word % 2 ^ 63 / 2 ^ 52 = 2047 then 2147483647 else _)
    split
    · rfl
    · by_cases subnormal : word % 2 ^ 63 / 2 ^ 52 = 0
      · have eq : word % 2 ^ 52 = word % 2 ^ 63 := by omega
        simp only [subnormal, ite_true, eq, Int.natCast_add, Int.natCast_one]
        omega
      · simp [subnormal]

end GoLean.Math
