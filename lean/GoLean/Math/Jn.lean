import GoLean.Math.J0
import GoLean.Math.J1
import GoLean.Math.Log

namespace GoLean.Math

/-- All integer indices in the bounded Bessel models convert exactly. -/
def besselIndexWord (n : Nat) : Nat := roundUnits (n * 2 ^ 1074) 1

/-- Sine and cosine factors for the leading large-argument expansion. -/
def besselHugePhase (arm64 fused : Bool) (n : Nat) (word : BitVec 64) :
    Option (Nat × Nat) := do
  let (sine, cosine) ← sincosWords arm64 fused word
  some (match n % 4 with
    | 0 => (subtractWords sine cosine, addWords cosine sine)
    | 1 => (subtractWords (negateWord sine) cosine, addWords (negateWord cosine) sine)
    | 2 => (addWords (negateWord sine) cosine, subtractWords (negateWord cosine) sine)
    | _ => (addWords sine cosine, subtractWords cosine sine))

private def jnForward (arm64 fused : Bool) (n : Nat) (word : BitVec 64) : Option Nat := do
  let j0 ← j0Word arm64 fused word
  let j1 ← j1Word arm64 fused word
  let (_, result) := (List.range (n - 1)).foldl (fun (a, b) i =>
    let ratio := divideWords (besselIndexWord (2 * (i + 1))) word.toNat
    (b, multiplyAddWords arm64 b ratio (negateWord a))) (j0, j1)
  some result

private def jnTaylor (n x : Nat) : Nat :=
  if 33 < n then 0
  else
    let half := multiplyWords x 0x3fe0000000000000
    let (factorial, power) := (List.range (n - 1)).foldl (fun (a, b) i =>
      (multiplyWords a (besselIndexWord (i + 2)), multiplyWords b half))
      (0x3ff0000000000000, half)
    divideWords power factorial

private def jnTerms (arm64 : Bool) (step : Nat) :
    Nat → Nat → Nat → Nat → Nat → Option Nat
  | 0, k, _, _, current => if current < 0x41cdcd6500000000 then none else some k
  | fuel + 1, k, z, previous, current =>
    if current < 0x41cdcd6500000000 then
      let next := addWords z step
      jnTerms arm64 step fuel (k + 1) next current
        (multiplyAddWords arm64 next current (negateWord previous))
    else some k

private def jnBackward (arm64 fused : Bool) (n : Nat) (word : BitVec 64) : Option (Option Nat) := do
  let one := 0x3ff0000000000000
  let x := word.toNat
  let w := divideWords (besselIndexWord (2 * n)) x
  let h := divideWords 0x4000000000000000 x
  let z := addWords w h
  let q1 := multiplyAddWords arm64 w z (negateWord one)
  let k ← jnTerms arm64 h 128 1 z w q1
  let initial := (List.range (k + 1)).foldl (fun t i =>
    divideWords one (subtractWords (divideWords (besselIndexWord (2 * (n + k - i))) x) t)) 0
  let order := besselIndexWord n
  let estimate ← logWord arm64 (BitVec.ofNat 64 (multiplyWords h order))
  let scale := 0x40862e42fefa39ef ≤ multiplyWords order estimate
  let (_, b, t) := (List.range (n - 1)).foldl (fun (a, b, t) i =>
    let factor := besselIndexWord (2 * (n - 1 - i))
    let next := subtractWords (divideWords (multiplyWords b factor) x) a
    if scale && 0x54b249ad2594c37d < next then
      (divideWords b next, one, divideWords t next)
    else (b, next, t)) (initial, one, initial)
  let j0 ← j0Word arm64 fused word
  let numerator := multiplyWords t j0
  some (if b % 2 ^ 63 = 0 then
    let infinity := 0x7ff0000000000000
    if numerator % 2 ^ 63 = 0 then none
    else some (if decide (2 ^ 63 ≤ numerator) != decide (2 ^ 63 ≤ b) then negateWord infinity else infinity)
  else some (divideWords numerator b))

/-- Go 1.27.0 src/math/jn.go:54–224, for signed orders -64 through 64.
    The outer option distinguishes supported evaluation from an out-of-scope
    order or exhausted recurrence guard. The inner option represents NaN. -/
def jnWord (arm64 fused : Bool) (order : Int) (word : BitVec 64) : Option (Option Nat) :=
  let n := order.natAbs
  let x := word.toNat % 2 ^ 63
  if 64 < n then none
  else if 0x7ff0000000000000 < x then some none
  else if x = 0x7ff0000000000000 then some (some 0)
  else if n = 0 then some (j0Word arm64 fused word)
  else if x = 0 then some (some 0)
  else if n = 1 then some (j1Word arm64 fused
    (if order < 0 then BitVec.ofNat 64 (negateWord word.toNat) else word))
  else
    let magnitude := BitVec.ofNat 64 x
    let result := if besselIndexWord n ≤ x then
      some (if 0x52d0000000000000 ≤ x then (do
        let (_, phase) ← besselHugePhase arm64 fused n magnitude
        let root ← sqrtWord magnitude
        some (divideWords (multiplyWords 0x3fe20dd750429b6d phase) root))
      else jnForward arm64 fused n magnitude)
    else if x < 0x3e20000000000000 then some (some (jnTaylor n x))
    else jnBackward arm64 fused n magnitude
    result.map (Option.map fun b =>
      let negative := decide (2 ^ 63 ≤ word.toNat) != decide (order < 0)
      if negative && n % 2 = 1 then negateWord b else b)

theorem jnWord_zeroOrder (arm64 fused : Bool) (word : BitVec 64)
    (finite : word.toNat % 2 ^ 63 < 0x7ff0000000000000) :
    jnWord arm64 fused 0 word = some (j0Word arm64 fused word) := by
  simp [jnWord, show ¬0x7ff0000000000000 < word.toNat % 2 ^ 63 by omega,
    show word.toNat % 2 ^ 63 ≠ 0x7ff0000000000000 by omega]

theorem jnWord_infinity (arm64 fused : Bool) (order : Int) (word : BitVec 64)
    (bounded : order.natAbs ≤ 64) (infinite : word.toNat % 2 ^ 63 = 0x7ff0000000000000) :
    jnWord arm64 fused order word = some (some 0) := by
  simp [jnWord, infinite, show ¬64 < order.natAbs by omega]

end GoLean.Math
