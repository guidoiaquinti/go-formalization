import GoLean.Math.Jn
import GoLean.Math.Y0
import GoLean.Math.Y1

namespace GoLean.Math

private def ynStep (arm64 : Bool) (x : Nat) (state : Nat × Nat) (i : Nat) : Nat × Nat :=
  let (a, b) := state
  let positiveInfinity := 0x7ff0000000000000
  let negativeInfinity := 0xfff0000000000000
  if b = negativeInfinity then state
  else
    let factor := divideWords (besselIndexWord (2 * (i + 1))) x
    let next := if factor = positiveInfinity then negativeInfinity
      else if arm64 then multiplyAddWords true factor b (negateWord a)
      else
        let product := multiplyWords factor b
        if product = negativeInfinity then product else subtractWords product a
    (b, next)

private def ynPositive (arm64 fused : Bool) (n : Nat) (word : BitVec 64) : Option Nat := do
  if 0x52d0000000000000 ≤ word.toNat then
    let (phase, _) ← besselHugePhase arm64 fused n word
    let root ← sqrtWord word
    some (divideWords (multiplyWords 0x3fe20dd750429b6d phase) root)
  else
    let y0 ← y0Word arm64 fused word
    let y1 ← y1Word arm64 fused word
    let (_, result) := (List.range (n - 1)).foldl (ynStep arm64 word.toNat) (y0, y1)
    some result

/-- Go 1.27.0 src/math/jn.go:235–306, for signed orders -64 through 64.
    None denotes NaN or an out-of-scope order. -/
def ynWord (arm64 fused : Bool) (order : Int) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  let magnitude := x % 2 ^ 63
  let n := order.natAbs
  let reverseSign := order < 0 && n % 2 = 1
  if 64 < n then none
  else if magnitude ≠ 0 && (2 ^ 63 ≤ x ∨ 0x7ff0000000000000 < magnitude) then none
  else if x = 0x7ff0000000000000 then some 0
  else if n = 0 then y0Word arm64 fused word
  else if magnitude = 0 then
    some (if reverseSign then 0x7ff0000000000000 else 0xfff0000000000000)
  else
    let result := if n = 1 then y1Word arm64 fused word else ynPositive arm64 fused n word
    result.map fun b => if reverseSign then negateWord b else b

theorem ynWord_zeroSignedOrder (arm64 fused : Bool) (order : Int) (word : BitVec 64)
    (bounded : order.natAbs ≤ 64) (nonzeroOrder : order.natAbs ≠ 0)
    (zero : word.toNat % 2 ^ 63 = 0) :
    ynWord arm64 fused order word =
      some (if order < 0 && order.natAbs % 2 = 1 then 0x7ff0000000000000 else 0xfff0000000000000) := by
  simp [ynWord, zero, nonzeroOrder, show ¬64 < order.natAbs by omega,
    show word.toNat ≠ 0x7ff0000000000000 by omega]

theorem ynWord_infinity (arm64 fused : Bool) (order : Int) (bounded : order.natAbs ≤ 64) :
    ynWord arm64 fused order (BitVec.ofNat 64 0x7ff0000000000000) = some 0 := by
  simp [ynWord, show ¬64 < order.natAbs by omega]

end GoLean.Math
