import GoLean.Math.Trunc
import GoLean.Math.Frexp

namespace GoLean.Math

/-- Encode the exact nonzero residue after subtracting the truncated
    integer part. The residue is counted in units of the input's last bit. -/
private def modfResidue (sign exponent residue : Nat) : Nat :=
  let top := residue.log2
  let payload := (residue <<< (52 - top)) % 2 ^ 52
  sign + (exponent + top - 52) * 2 ^ 52 + payload

/-- Go computes the integer part with Trunc, subtracts it, and copies the
    original sign. For a finite input the subtraction is exact: the
    operands are within a factor of two or the integer part is zero.
    None represents the fractional NaN for an infinite or NaN input.
    https://github.com/golang/go/blob/go1.27.0/src/math/modf.go#L7-L18
    https://github.com/golang/go/blob/go1.27.0/src/math/floor.go#L61-L80 -/
def modfWords (word : BitVec 64) : Nat × Option Nat :=
  let bits := word.toNat
  let whole := trunc word
  let exponent := bits / 2 ^ 52 % 2048
  let sign := bits - bits % 2 ^ 63
  let fraction :=
    if exponent = 2047 then none
    else if exponent < 1023 then some bits
    else
      let residue := bits - whole
      if residue = 0 then some sign
      else some (modfResidue sign exponent residue)
  (whole, fraction)

/-- The integer result clears all bits below the unit position. The
    fractional result normalizes the corresponding remainder and
    retains the input's sign, including when either result is zero. -/
theorem modfWords_spec (word : BitVec 64) :
    modfWords word =
      let bits := word.toNat
      let exponent := bits / 2 ^ 52 % 2048
      let sign := bits - bits % 2 ^ 63
      if exponent = 2047 then (bits, none)
      else if exponent < 1023 then (sign, some bits)
      else if exponent < 1075 then
        let unit := 2 ^ (1075 - exponent)
        let residue := bits % unit
        let fraction :=
          if residue = 0 then sign
          else
            let top := residue.log2
            sign + (exponent + top - 52) * 2 ^ 52 +
              (residue - 2 ^ top) * 2 ^ (52 - top)
        (bits / unit * unit, some fraction)
      else (bits, some sign) := by
  unfold modfWords
  rw [trunc_spec]
  dsimp only
  let bits := word.toNat
  let exponent := bits / 2 ^ 52 % 2048
  by_cases special : exponent = 2047
  · rw [ite_eq_left special, ite_eq_left special,
      ite_eq_right (by omega : ¬ exponent < 1023), ite_eq_right (by omega : ¬ exponent < 1075)]
  · rw [ite_eq_right special, ite_eq_right special]
    by_cases belowOne : exponent < 1023
    · rw [ite_eq_left belowOne, ite_eq_left belowOne, ite_eq_left belowOne]
    · rw [ite_eq_right belowOne, ite_eq_right belowOne, ite_eq_right belowOne]
      by_cases fractional : exponent < 1075
      · rw [ite_eq_left fractional, ite_eq_left fractional]
        let unit := 2 ^ (1075 - exponent)
        have residue : bits - bits / unit * unit = bits % unit := by
          have parts := Nat.mod_add_div' bits unit
          omega
        rw [residue]
        by_cases empty : bits % unit = 0
        · rw [ite_eq_left empty, ite_eq_left empty]
        · rw [ite_eq_right empty, ite_eq_right empty]
          have residueBound : bits % unit < 2 ^ 52 :=
            Nat.lt_of_lt_of_le (Nat.mod_lt _ (Nat.two_pow_pos _))
              (Nat.pow_le_pow_right Nat.zero_lt_two (by omega))
          rw [modfResidue, normalizedPayload _ empty residueBound, Nat.shiftLeft_eq]
      · rw [ite_eq_right fractional, ite_eq_right fractional, Nat.sub_self, ite_eq_left rfl]

end GoLean.Math
