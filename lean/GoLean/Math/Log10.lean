import GoLean.Math.Log

namespace GoLean.Math

/-- Go 1.27.0 src/math/log10.go:9–18. The generic wrapper multiplies
    the already rounded public Log result by the rounded reciprocal
    of ln(10). On amd64 that call uses the architecture Log assembly,
    including its treatment of subnormal inputs. -/
def log10Word (arm64 : Bool) (word : BitVec 64) : Option Nat :=
  (logWord arm64 word).map fun logarithm =>
    if 0x7ff0000000000000 ≤ logarithm % 2 ^ 63 then logarithm
    else multiplyWords logarithm 0x3fdbcb7b1526e50e

theorem log10Word_zero (arm64 : Bool) (word : BitVec 64)
    (zero : word.toNat % 2 ^ 63 = 0) :
    log10Word arm64 word = some 0xfff0000000000000 := by
  rw [log10Word, logWord_zero _ _ zero]
  rfl

theorem log10Word_negative (arm64 : Bool) (word : BitVec 64)
    (negative : 2 ^ 63 ≤ word.toNat) (nonzero : word.toNat % 2 ^ 63 ≠ 0) :
    log10Word arm64 word = none := by
  rw [log10Word, logWord_negative _ _ negative nonzero]
  rfl

theorem log10Word_inf (arm64 : Bool) :
    log10Word arm64 (BitVec.ofNat 64 0x7ff0000000000000) = some 0x7ff0000000000000 := by
  rw [log10Word, logWord_inf]
  rfl

theorem log10Word_scale (arm64 : Bool) (word : BitVec 64) (result : Nat)
    (logarithm : logWord arm64 word = some result)
    (finite : result % 2 ^ 63 < 0x7ff0000000000000) :
    log10Word arm64 word = some (multiplyWords result 0x3fdbcb7b1526e50e) := by
  simp [log10Word, logarithm, Nat.not_le.mpr finite]

end GoLean.Math
