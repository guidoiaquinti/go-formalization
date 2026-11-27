import GoLean.Math.BesselRational

namespace GoLean.Math

/-- Go 1.27.0 src/math/j0.go:410–429. Selects the source interval,
    evaluates the rounded ratio, subtracts 1/8, and divides by x. None denotes the
    nil-coefficient-table panic below two or on NaN. -/
def qzeroWord (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  if x < 0x4000000000000000 ∨ 0x7ff0000000000000 < x then none
  else if x = 0x7ff0000000000000 then some 0x8000000000000000
  else
    let (p, q) := besselCoefficients x
      ([
        0x40e212d40e901566, 0x40c14d993e18f46d, 0x40816d6315301825,
        0x402789525bb334d6, 0x3fb2bffffffffe2c, 0x0000000000000000], [
        0xc114fd6d2c9530c5, 0x4129a66b28de0b3d, 0x412883da83a52b43,
        0x4101665254d38c3f, 0x40bfa2584e6b0563, 0x406478d5365b39bc])
      ([
        0x409f17e953c6e3a6, 0x40900cf99dc8c481, 0x4060e3920a8788e9,
        0x401757b0b9953dd3, 0x3fb2bfffd172b04c, 0x3db43d8f29cc8cd9], [
        0xc0b4ea57bedbc609, 0x40e191181f7a54a0, 0x40ebb5e397e02372,
        0x40d267d27b591e6d, 0x40a03ba0da21c0ce, 0x4054b1b3fb5e1543])
      ([
        0x4064d77c81fa21e0, 0x406559dbe25efd1f, 0x40454f98962daedd,
        0x400ac0fc61149cf5, 0x3fb2bfee0e8d0842, 0x3e32cd036adecb82], [
        0xc062a7eb201cf40f, 0x40a3a8aad94fb1c0, 0x40b93c6cd7c76a28,
        0x40acf04be44dfc63, 0x40862d8386544eb3, 0x40486122bfe343a6])
      ([
        0x403040b171814bb4, 0x403faa8e29fbdc4a, 0x402cfdbfaaf96fe5,
        0x3ffff897e727779c, 0x3fb2bec53e883e34, 0x3e84313b54f76bdb], [
        0xc0153e6af8b32931, 0x406a95530e001365, 0x408b977c9c5cc214,
        0x408a664522b3bf22, 0x4070d591e4d14b40, 0x403e5d96f7c07aed])
    let ratio := besselRational arm64 fused x p q
    some (divideWords (addWords 0xbfc0000000000000 ratio) x)

theorem qzeroWord_panic (arm64 fused : Bool) (word : BitVec 64)
    (outside : word.toNat < 0x4000000000000000 ∨ 0x7ff0000000000000 < word.toNat) :
    qzeroWord arm64 fused word = none := by
  simp [qzeroWord, outside]

theorem qzeroWord_infinity (arm64 fused : Bool) :
    qzeroWord arm64 fused (BitVec.ofNat 64 0x7ff0000000000000) = some 0x8000000000000000 := by
  cases arm64 <;> cases fused <;> decide

end GoLean.Math
