import GoLean.Math.FmaSplit
import GoLean.Math.Hypot

namespace GoLean.Math

private def fmaInfinity : Nat := 0x7ff0000000000000

private def stickyShift (value shift : Nat) : Nat :=
  (value >>> shift) ||| (if value % 2 ^ shift = 0 then 0 else 1)

private def roundFma (sgn : Nat) (exponent : Int) (mantissa : Nat) : Nat :=
  if 2045 < exponent ∨ (exponent = 2045 ∧ (mantissa + 2 ^ 9) >>> 63 = 1) then
    sgn * 2 ^ 63 + fmaInfinity
  else
    let m := if exponent < 0 then stickyShift mantissa (-exponent).toNat else mantissa
    let rounded := (m + 2 ^ 9) >>> 10
    let even := if m % 2 ^ 10 = 2 ^ 9 then rounded - rounded % 2 else rounded
    sgn * 2 ^ 63 + (if even = 0 then 0 else exponent.toNat * 2 ^ 52) + even

private def alignFma (ps : Nat) (pe : Int) (product zs : Nat) (ze : Int) (addend : Nat) : Nat :=
  if ps ≠ zs ∧ pe = ze ∧ product = addend then 0
  else
    let aligned := stickyShift addend (pe - ze).toNat
    if ps = zs then
      let sum := product + aligned
      let carry := sum >>> 127
      let m := stickyShift sum (64 + carry) % 2 ^ 64
      roundFma ps (pe - 1 + (carry : Int)) m
    else
      let diff := product - aligned
      let lz := 127 - diff.log2
      let shifted := diff <<< (lz - 1)
      roundFma ps (pe - (lz : Int)) (stickyShift shifted 64)

/-- Go 1.27.0 src/math/fma.go:95–182. The pair of uint64 limbs is
    represented by their joined natural value. The model preserves the
    product normalization, sticky alignment, addition/cancellation,
    renormalization, underflow shift and ties-to-even rounding. None
    represents NaN by classification. -/
def fmaWords (x y z : BitVec 64) : Option Nat :=
  let bx := x.toNat % 2 ^ 63
  let ymag := y.toNat % 2 ^ 63
  let bz := z.toNat % 2 ^ 63
  let sx := x.toNat >>> 63
  let sy := y.toNat >>> 63
  let sz := z.toNat >>> 63
  let ps := sx ^^^ sy
  if fmaInfinity < bx ∨ fmaInfinity < ymag ∨ fmaInfinity < bz then none
  else if bx = fmaInfinity ∨ ymag = fmaInfinity then
    if bx = 0 ∨ ymag = 0 ∨ (bz = fmaInfinity ∧ ps ≠ sz) then none
    else some (ps * 2 ^ 63 + fmaInfinity)
  else if bx = 0 ∨ ymag = 0 then
    if bz = 0 then some ((ps &&& sz) * 2 ^ 63) else some z.toNat
  else if bz = 0 then
    some (ps * 2 ^ 63 + roundUnits (magnitudeUnits bx * magnitudeUnits ymag) (2 ^ 1074))
  else if bz = fmaInfinity then some z.toNat
  else
    let (_, xe, xm) := fmaSplit x
    let (_, ye, ym) := fmaSplit y
    let (_, ze, zm) := fmaSplit z
    let mantissa := (xm * ym) <<< 21
    let shift : Nat := if mantissa < 2 ^ 126 then 1 else 0
    let product := mantissa <<< shift
    let pe := xe + ye - 1022 - (shift : Int)
    let addend := zm <<< 74
    if pe < ze ∨ (pe = ze ∧ product >>> 64 < addend >>> 64) then
      some (alignFma sz ze addend ps pe product)
    else some (alignFma ps pe product sz ze addend)

theorem fmaWords_nan (x y z : BitVec 64)
    (h : 0x7ff0000000000000 < x.toNat % 2 ^ 63 ∨
      0x7ff0000000000000 < y.toNat % 2 ^ 63 ∨ 0x7ff0000000000000 < z.toNat % 2 ^ 63) :
    fmaWords x y z = none := by
  simp [fmaWords, fmaInfinity, h]

theorem fmaWords_zeroAddend (x y z : BitVec 64)
    (hx : 0 < x.toNat % 2 ^ 63 ∧ x.toNat % 2 ^ 63 < 0x7ff0000000000000)
    (hy : 0 < y.toNat % 2 ^ 63 ∧ y.toNat % 2 ^ 63 < 0x7ff0000000000000)
    (hz : z.toNat % 2 ^ 63 = 0) :
    fmaWords x y z = some (((x.toNat >>> 63) ^^^ (y.toNat >>> 63)) * 2 ^ 63 +
      roundUnits (magnitudeUnits (x.toNat % 2 ^ 63) * magnitudeUnits (y.toNat % 2 ^ 63)) (2 ^ 1074)) := by
  unfold fmaWords fmaInfinity
  dsimp only
  rw [ite_eq_right (by omega), ite_eq_right (by omega), ite_eq_right (by omega), ite_eq_left hz]

theorem fmaWords_cancel (x y z : BitVec 64)
    (hx : 0 < x.toNat % 2 ^ 63 ∧ x.toNat % 2 ^ 63 < 0x7ff0000000000000)
    (hy : 0 < y.toNat % 2 ^ 63 ∧ y.toNat % 2 ^ 63 < 0x7ff0000000000000)
    (hz : 0 < z.toNat % 2 ^ 63 ∧ z.toNat % 2 ^ 63 < 0x7ff0000000000000)
    (signs : (x.toNat >>> 63 ^^^ y.toNat >>> 63) ≠ z.toNat >>> 63)
    (cancel :
      let (_, xe, xm) := fmaSplit x
      let (_, ye, ym) := fmaSplit y
      let (_, ze, zm) := fmaSplit z
      let mantissa := (xm * ym) <<< 21
      let shift : Nat := if mantissa < 2 ^ 126 then 1 else 0
      xe + ye - 1022 - (shift : Int) = ze ∧ mantissa <<< shift = zm <<< 74) :
    fmaWords x y z = some 0 := by
  unfold fmaWords fmaInfinity
  dsimp only at cancel ⊢
  rw [ite_eq_right (by omega), ite_eq_right (by omega), ite_eq_right (by omega),
    ite_eq_right (by omega), ite_eq_right (by omega)]
  split <;> simp_all [alignFma]

end GoLean.Math
