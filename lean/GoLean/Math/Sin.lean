import GoLean.Math.Exp2

namespace GoLean.Math

private def pi4Digits : Array Nat := #[
  0x0000000000000001, 0x45f306dc9c882a53, 0xf84eafa3ea69bb81, 0xb6c52b3278872083,
  0xfca2c757bd778ac3, 0x6e48dc74849ba5c0, 0x0c925dd413a32439, 0xfc3bd63962534e7d,
  0xd1046bea5d768909, 0xd338e04d68befc82, 0x7323ac7306a673e9, 0x3908bf177bf25076,
  0x3ff12fffbc0b301f, 0xde5e2316b414da3e, 0xda6cfd9e4f96136e, 0x9e8c7ecd3cbfd45a,
  0xea4f758fd7cbe2f6, 0x7a0e73ef14a525d4, 0xd7f6bf623f1aba10, 0xac06608df8f6d757]

private def limb (word : Nat) : Nat := word % 2 ^ 64

private def shiftedPi4Digit (index shift : Nat) : Nat :=
  limb (pi4Digits[index]! <<< shift) ||| (pi4Digits[index + 1]! >>> (64 - shift))

def reduceLargeSin (magnitude : Nat) : Nat × Nat :=
  let exponent := (magnitude >>> 52) - 1014
  let digit := exponent / 64
  let shift := exponent % 64
  let mantissa := (magnitude % 2 ^ 52) + 2 ^ 52
  let z0 := shiftedPi4Digit digit shift
  let z1 := shiftedPi4Digit (digit + 1) shift
  let z2 := shiftedPi4Digit (digit + 2) shift
  let low := limb (z1 * mantissa) + ((z2 * mantissa) >>> 64)
  let high := limb (z0 * mantissa + ((z1 * mantissa) >>> 64) + (low >>> 64))
  let low := limb low
  let octant := high >>> 61
  let high := limb (high <<< 3) ||| (low >>> 61)
  let leadingZeros := if high = 0 then 64 else 63 - high.log2
  let tail := if leadingZeros < 64 then low >>> (63 - leadingZeros) else 0
  let normalized := limb (high <<< (leadingZeros + 1)) ||| tail
  let encoded := (normalized >>> 12) ||| ((1022 - leadingZeros) <<< 52)
  let fraction := if octant % 2 = 1 then subtractWords encoded 0x3ff0000000000000 else encoded
  ((octant + octant % 2) % 8, multiplyWords fraction 0x3fe921fb54442d18)

def reduceSmallSin (arm64 : Bool) (magnitude : Nat) : Nat × Nat :=
  let scaled := multiplyWords magnitude 0x3ff45f306dc9c883
  let truncated := magnitudeUnits scaled / 2 ^ 1074
  let even : Nat := truncated + truncated % 2
  let count := roundSigned ((even : Int) * 2 ^ 1074) 1
  let a := multiplyAddWords arm64 count 0xbfe921fb40000000 magnitude
  let b := multiplyAddWords arm64 count 0xbe64442d00000000 a
  let c := multiplyAddWords arm64 count 0xbce8469898cc5170 b
  (even % 8, c)

def sinPolynomial (arm64 fused : Bool) (z : Nat) : Nat :=
  let zz := multiplyWords z z
  let term := multiplyAddWords arm64 0x3de5d8fd1fd19ccd zz 0xbe5ae5e5a9291f5d
  let term := multiplyAddWords arm64 term zz 0x3ec71de3567d48a1
  let term := multiplyAddWords arm64 term zz 0xbf2a01a019bfdf03
  let term := multiplyAddWords arm64 term zz 0x3f8111111110f7d0
  let term := multiplyAddWords arm64 term zz 0xbfc5555555555548
  multiplyAddWords (arm64 || fused) (multiplyWords z zz) term z

def cosPolynomial (arm64 fused : Bool) (z : Nat) : Nat :=
  let zz := multiplyWords z z
  let term := multiplyAddWords arm64 0xbda8fa49a0861a9b zz 0x3e21ee9d7b4e3f05
  let term := multiplyAddWords arm64 term zz 0xbe927e4f7eac4bc6
  let term := multiplyAddWords arm64 term zz 0x3efa01a019c844f5
  let term := multiplyAddWords arm64 term zz 0xbf56c16c16c14f91
  let term := multiplyAddWords arm64 term zz 0x3fa555555555554b
  let linear := multiplyAddWords arm64 zz 0xbfe0000000000000 0x3ff0000000000000
  multiplyAddWords (arm64 || fused) (multiplyWords zz zz) term linear

/-- Go 1.27.0 src/math/sin.go:185–244 and src/math/trig_reduce.go:31–102.
    Preserve both reductions, octant reflection, the two coefficient chains,
    and compiler-contracted rounding. The large path retains the three
    selected digits of 4/π and the fixed-width limb operations. -/
def sinWord (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  let raw := word.toNat
  let magnitude := raw % 2 ^ 63
  if magnitude = 0 then some raw
  else if 0x7ff0000000000000 ≤ magnitude then none
  else
    let (octant, reduced) := if magnitude < 0x41c0000000000000 then reduceSmallSin arm64 magnitude
      else reduceLargeSin magnitude
    let result := if octant % 4 = 2 then cosPolynomial arm64 fused reduced else sinPolynomial arm64 fused reduced
    some (if decide (2 ^ 63 ≤ raw) != decide (4 ≤ octant) then negateWord result else result)

theorem sinWord_zero (arm64 fused : Bool) (word : BitVec 64)
    (zero : word.toNat % 2 ^ 63 = 0) :
    sinWord arm64 fused word = some word.toNat := by
  simp [sinWord, zero]

theorem sinWord_nonfinite (arm64 fused : Bool) (word : BitVec 64)
    (nonfinite : 0x7ff0000000000000 ≤ word.toNat % 2 ^ 63) :
    sinWord arm64 fused word = none := by
  simp [sinWord, nonfinite, show word.toNat % 2 ^ 63 ≠ 0 by omega]

theorem sinWord_opposite (arm64 fused : Bool) (x y : BitVec 64)
    (positive : x.toNat < 2 ^ 63) (negative : 2 ^ 63 ≤ y.toNat)
    (same : x.toNat % 2 ^ 63 = y.toNat % 2 ^ 63)
    (finite : x.toNat % 2 ^ 63 < 0x7ff0000000000000)
    (nonzero : x.toNat % 2 ^ 63 ≠ 0) :
    sinWord arm64 fused y = (sinWord arm64 fused x).map negateWord := by
  unfold sinWord
  simp only [← same, nonzero, Nat.not_le.mpr finite, ite_false]
  split <;> rename_i range
  all_goals
    simp only [negative, Nat.not_le.mpr positive, decide_true, decide_false,
      Option.map_some]
    split
    next firstHalf =>
      simp at firstHalf
      have secondHalf := Nat.not_le.mpr firstHalf
      simp [negateWord, secondHalf]
    next =>
      simp_all [negateWord, Nat.xor_assoc]

end GoLean.Math
