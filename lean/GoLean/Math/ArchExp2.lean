import GoLean.Math.Exp2

namespace GoLean.Math

/-- Go 1.27.0 src/math/exp2_noasm.go:11–13 and exp_arm64.s:112–182.
    The arm64 body evaluates the exponential. Ports selected by the Go
    fallback build constraint have an unconditional panic stub. -/
def archExp2Word (arm64 : Bool) (word : BitVec 64) : Except String Nat :=
  if arm64 then .ok (exp2Word true false false word) else .error "not implemented"

theorem archExp2Word_stub (word : BitVec 64) :
    archExp2Word false word = .error "not implemented" := rfl

theorem archExp2Word_nan (word : BitVec 64)
    (nan : 0x7ff0000000000000 < word.toNat % 2 ^ 63) :
    archExp2Word true word = .ok word.toNat := by
  simp [archExp2Word, exp2Word_nan true false false word nan]

theorem archExp2Word_zeros :
    archExp2Word true (0 : BitVec 64) = .ok 0x3ff0000000000000 ∧
      archExp2Word true (BitVec.ofNat 64 (2 ^ 63)) = .ok 0x3ff0000000000000 := by
  rcases exp2Word_zeros true false false with ⟨positive, negative⟩
  exact ⟨congrArg Except.ok positive, congrArg Except.ok negative⟩

end GoLean.Math
