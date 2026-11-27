import GoLean.Math.Log

namespace GoLean.Math

/-- Go 1.27.0 src/math/log_asm.go:11, log_amd64.s:22–112 and
    log_stub.go:11–13. The selected amd64 assembly supplies a leading
    significand bit even for positive subnormals. The fallback panics. -/
def archLogWord (amd64 : Bool) (word : BitVec 64) : Except String (Option Nat) :=
  if amd64 then .ok (logWord false word) else .error "not implemented"

theorem archLogWord_stub (word : BitVec 64) :
    archLogWord false word = .error "not implemented" := rfl

theorem archLogWord_zero (word : BitVec 64) (zero : word.toNat % 2 ^ 63 = 0) :
    archLogWord true word = .ok (some 0xfff0000000000000) :=
  congrArg Except.ok (logWord_zero false word zero)

theorem archLogWord_negative (word : BitVec 64)
    (negative : 2 ^ 63 ≤ word.toNat) (nonzero : word.toNat % 2 ^ 63 ≠ 0) :
    archLogWord true word = .ok none :=
  congrArg Except.ok (logWord_negative false word negative nonzero)

end GoLean.Math
