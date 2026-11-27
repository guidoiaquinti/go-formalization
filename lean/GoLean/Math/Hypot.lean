import GoLean.Math.Mod
import GoLean.Math.Sqrt

namespace GoLean.Math

private def positiveInfinity : Nat := 0x7ff0000000000000

/-- Round a nonnegative fraction of the smallest binary64 unit to an
    encoded binary64 word using nearest, ties to even. -/
def roundUnits (numerator denominator : Nat) : Nat :=
  let shift := (numerator / denominator).log2 - 52
  let scale := 2 ^ shift
  let divisor := denominator * scale
  let whole := numerator / divisor
  let residue := numerator % divisor
  let rounded := if divisor < 2 * residue ∨ (divisor = 2 * residue ∧ whole % 2 = 1)
    then whole + 1 else whole
  min positiveInfinity (encodeUnits (rounded * scale))

private def hypotRadicand (fused : Bool) (small large : Nat) : Nat :=
  let unitOne := 2 ^ 1074
  let ratio := magnitudeUnits (roundUnits (small * unitOne) large)
  let square := ratio * ratio
  if fused then roundUnits (unitOne * unitOne + square) unitOne
  else roundUnits (unitOne + magnitudeUnits (roundUnits square unitOne)) 1

/-- Go 1.27.0 src/math/hypot.go:20–44 and hypot_amd64.s:13–50.
    Strip signs, handle exceptional operands, order the magnitudes, and
    round the divide, square/add, root and multiply in source order.
    The fused profile represents the Go arm64 compiler output.
    None denotes NaN by classification. -/
def hypotWords (fused : Bool) (x y : BitVec 64) : Option Nat :=
  let p := x.toNat % 2 ^ 63
  let q := y.toNat % 2 ^ 63
  if p = positiveInfinity ∨ q = positiveInfinity then some positiveInfinity
  else if positiveInfinity < p ∨ positiveInfinity < q then none
  else
    let large := magnitudeUnits (max p q)
    if large = 0 then some 0
    else
      let unitOne := 2 ^ 1074
      let radicand := hypotRadicand fused (magnitudeUnits (min p q)) large
      (sqrtWord (BitVec.ofNat 64 radicand)).map fun root =>
        roundUnits (large * magnitudeUnits root) unitOne

theorem hypotWords_symm (fused : Bool) (x y : BitVec 64) :
    hypotWords fused x y = hypotWords fused y x := by
  simp only [hypotWords, or_comm, Nat.max_comm, Nat.min_comm]

theorem hypotWords_special (fused : Bool) (x y : BitVec 64)
    (h : 0x7ff0000000000000 ≤ x.toNat % 2 ^ 63 ∨
      0x7ff0000000000000 ≤ y.toNat % 2 ^ 63) :
    hypotWords fused x y =
      if x.toNat % 2 ^ 63 = 0x7ff0000000000000 ∨ y.toNat % 2 ^ 63 = 0x7ff0000000000000
      then some 0x7ff0000000000000 else none := by
  unfold hypotWords positiveInfinity
  dsimp only
  by_cases inf : x.toNat % 2 ^ 63 = 0x7ff0000000000000 ∨ y.toNat % 2 ^ 63 = 0x7ff0000000000000
  · rw [ite_eq_left inf, ite_eq_left inf]
  · have nan : 0x7ff0000000000000 < x.toNat % 2 ^ 63 ∨ 0x7ff0000000000000 < y.toNat % 2 ^ 63 := by omega
    rw [ite_eq_right inf, ite_eq_left nan, ite_eq_right inf]

theorem hypotWords_zeros (fused : Bool) (x y : BitVec 64)
    (hx : x.toNat % 2 ^ 63 = 0) (hy : y.toNat % 2 ^ 63 = 0) :
    hypotWords fused x y = some 0 := by
  simp [hypotWords, positiveInfinity, hx, hy, magnitudeUnits]

end GoLean.Math
