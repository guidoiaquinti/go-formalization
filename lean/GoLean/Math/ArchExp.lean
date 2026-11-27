import GoLean.Math.Exp

namespace GoLean.Math

/-- Go 1.27.0 src/math/exp_asm.go:11, exp_arm64.s:26–105 and
    exp_amd64.s:40–157. Direct entry to the pinned assembly body. -/
def archExpWord (arm64 fused : Bool) (word : BitVec 64) : Nat :=
  expWord arm64 fused word

theorem archExpWord_nan (arm64 fused : Bool) (word : BitVec 64)
    (nan : 0x7ff0000000000000 < word.toNat % 2 ^ 63) :
    archExpWord arm64 fused word = word.toNat :=
  expWord_nan arm64 fused word nan

theorem archExpWord_infinities (arm64 fused : Bool) :
    archExpWord arm64 fused (BitVec.ofNat 64 0x7ff0000000000000) = 0x7ff0000000000000 ∧
    archExpWord arm64 fused (BitVec.ofNat 64 0xfff0000000000000) = 0 :=
  expWord_infinities arm64 fused

end GoLean.Math
