import GoLean.Math.Ilogb

namespace GoLean.Math

/-- A finite result represents the exact float64 conversion of its integral
    exponent. Special results carry their raw binary64 word. -/
inductive LogbResult where
  | finite (exponent : Int)
  | special (word : Nat)
  deriving DecidableEq

/-- Go checks signed zero, infinity, and NaN before converting ilogb.
    https://github.com/golang/go/blob/go1.27.0/src/math/logb.go#L7-L25
    https://github.com/golang/go/blob/go1.27.0/src/math/logb.go#L47-L52 -/
def logb (word : Nat) : LogbResult :=
  let e := (word >>> 52) &&& 2047
  let f := word &&& (2 ^ 52 - 1)
  if e = 0 ∧ f = 0 then .special 0xfff0000000000000
  else if e = 2047 ∧ f = 0 then .special 0x7ff0000000000000
  else if e = 2047 then .special word
  else .finite (ilogb word)

theorem logb_spec (word : Nat) :
    logb word =
      let mag := word % 2 ^ 63
      if mag = 0 then .special 0xfff0000000000000
      else if mag = 0x7ff0000000000000 then .special 0x7ff0000000000000
      else if 0x7ff0000000000000 < mag then .special word
      else if mag < 2 ^ 52 then .finite ((mag.log2 : Int) - 1074)
      else .finite (((mag / 2 ^ 52 : Nat) : Int) - 1023) := by
  unfold logb
  dsimp only
  rw [Nat.shiftRight_eq_div_pow, show (2047 : Nat) = 2 ^ 11 - 1 by decide,
    Nat.and_two_pow_sub_one_eq_mod, Nat.and_two_pow_sub_one_eq_mod]
  have exp : word / 2 ^ 52 % 2 ^ 11 = word % 2 ^ 63 / 2 ^ 52 := by omega
  have zero : (word % 2 ^ 63 / 2 ^ 52 = 0 ∧ word % 2 ^ 52 = 0) ↔
      word % 2 ^ 63 = 0 := by omega
  have inf : (word % 2 ^ 63 / 2 ^ 52 = 2047 ∧ word % 2 ^ 52 = 0) ↔
      word % 2 ^ 63 = 0x7ff0000000000000 := by omega
  rw [exp]
  by_cases z : word % 2 ^ 63 = 0
  · rw [ite_eq_left (zero.mpr z), ite_eq_left z]
  · rw [ite_eq_right (mt zero.mp z), ite_eq_right z]
    by_cases i : word % 2 ^ 63 = 0x7ff0000000000000
    · rw [ite_eq_left (inf.mpr i), ite_eq_left i]
    · rw [ite_eq_right (mt inf.mp i), ite_eq_right i]
      by_cases n : 0x7ff0000000000000 < word % 2 ^ 63
      · rw [ite_eq_left (by omega : word % 2 ^ 63 / 2 ^ 52 = 2047), ite_eq_left n]
      · rw [ite_eq_right (by omega : word % 2 ^ 63 / 2 ^ 52 ≠ 2047), ite_eq_right n, ilogb_spec,
          ite_eq_right z, ite_eq_right (by omega : word % 2 ^ 63 / 2 ^ 52 ≠ 2047)]
        by_cases small : word % 2 ^ 63 < 2 ^ 52
        · rw [ite_eq_left (by omega : word % 2 ^ 63 / 2 ^ 52 = 0), ite_eq_left small]
        · rw [ite_eq_right (by omega : word % 2 ^ 63 / 2 ^ 52 ≠ 0), ite_eq_right small]

end GoLean.Math
