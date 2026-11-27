import Init.Data.BitVec.Lemmas

namespace GoLean.Math

/-- On binary64 encodings, clearing the sign is reduction modulo 2^63.
    Negation toggles the sign by adding 2^63 modulo 2^64. The final ordered
    comparison exceeds MaxFloat64 only through the upper, non-NaN endpoint
    of the nonnegative encodings.

    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L42-L57
    https://github.com/golang/go/blob/go1.27.0/src/math/abs.go#L7-L15
    https://github.com/golang/go/blob/go1.27.0/src/math/const.go#L34 -/
def isInf (repr : BitVec 64) (sign : Int) : Bool :=
  let bits := repr.toNat
  let adjusted :=
    if sign = 0 then bits % 2 ^ 63
    else if sign < 0 then (bits + 2 ^ 63) % 2 ^ 64
    else bits
  decide (0x7fefffffffffffff < adjusted ∧ adjusted ≤ 0x7ff0000000000000)

theorem isInf_spec (repr : BitVec 64) (sign : Int) :
    isInf repr sign = true ↔
      (0 ≤ sign ∧ repr.toNat = 0x7ff0000000000000) ∨
      (sign ≤ 0 ∧ repr.toNat = 0xfff0000000000000) := by
  have bound := repr.isLt
  by_cases hz : sign = 0
  · simp [isInf, hz]
    omega
  by_cases hn : sign < 0
  · simp [isInf, hz, hn]
    omega
  · simp [isInf, hz, hn]
    omega

end GoLean.Math
