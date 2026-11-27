import GoLean.Math.Hypot

namespace GoLean.Math

/-- Go 1.27.0 src/math/hypot_asm.go:11, hypot_amd64.s:12–50 and
    hypot_noasm.go:11–13. The amd64 SSE2 routine uses split arithmetic
    and a canonical NaN word. The no-assembly ports panic. -/
def archHypotWord (amd64 : Bool) (x y : BitVec 64) : Except String Nat :=
  if amd64 then .ok ((hypotWords false x y).getD 0x7ff8000000000001)
  else .error "not implemented"

theorem archHypotWord_stub (x y : BitVec 64) :
    archHypotWord false x y = .error "not implemented" := rfl

theorem archHypotWord_symm (amd64 : Bool) (x y : BitVec 64) :
    archHypotWord amd64 x y = archHypotWord amd64 y x := by
  simp [archHypotWord, hypotWords_symm false x y]

theorem archHypotWord_special (x y : BitVec 64)
    (special : 0x7ff0000000000000 ≤ x.toNat % 2 ^ 63 ∨
      0x7ff0000000000000 ≤ y.toNat % 2 ^ 63) :
    archHypotWord true x y = .ok
      (if x.toNat % 2 ^ 63 = 0x7ff0000000000000 ∨ y.toNat % 2 ^ 63 = 0x7ff0000000000000
        then 0x7ff0000000000000 else 0x7ff8000000000001) := by
  unfold archHypotWord
  simp only [ite_true, hypotWords_special false x y special]
  split <;> rfl

end GoLean.Math
