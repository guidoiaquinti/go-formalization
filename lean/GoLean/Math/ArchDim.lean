import GoLean.Math.Max
import GoLean.Math.Min
import Init.Data.BitVec.Basic

namespace GoLean.Math

inductive ArchDimProfile where
  | amd64
  | arm64
  | fallback

private def dimNaN (word : Nat) : Bool :=
  decide (0x7ff0000000000000 < word % 2 ^ 63)

private def dimSignalingNaN (word : Nat) : Bool :=
  dimNaN word && (word &&& 0x0008000000000000) == 0

/-- FMAXD and FMIND receive y as their first operand. They prioritize
    signaling NaNs, then quiet NaNs, preserving the selected payload. -/
private def dimArmNaN (x y : Nat) : Option Nat :=
  if dimSignalingNaN y then some (y ||| 0x0008000000000000)
  else if dimSignalingNaN x then some (x ||| 0x0008000000000000)
  else if dimNaN y then some y
  else if dimNaN x then some x
  else none

private def dimAssembly (profile : ArchDimProfile) (dominant : Nat)
    (ordinary : Nat → Nat → Nat) (x y : Nat) : Except String Nat :=
  match profile with
  | .fallback => .error "not implemented"
  | .amd64 => .ok (ordinary x y)
  | .arm64 =>
    if x = dominant ∨ y = dominant then .ok dominant
    else .ok ((dimArmNaN x y).getD (ordinary x y))

/-- Go 1.27.0 src/math/dim_asm.go:11, dim_amd64.s:12–49,
    dim_arm64.s:12–29 and dim_noasm.go:11–13. -/
def archMaxWord (profile : ArchDimProfile) (x y : BitVec 64) : Except String Nat :=
  dimAssembly profile 0x7ff0000000000000 maxFloat64 x.toNat y.toNat

/-- Go 1.27.0 src/math/dim_asm.go:15, dim_amd64.s:52–89,
    dim_arm64.s:32–49 and dim_noasm.go:17–19. -/
def archMinWord (profile : ArchDimProfile) (x y : BitVec 64) : Except String Nat :=
  dimAssembly profile 0xfff0000000000000 minFloat64 x.toNat y.toNat

theorem archMaxWord_stub (x y : BitVec 64) :
    archMaxWord .fallback x y = .error "not implemented" := rfl

theorem archMinWord_stub (x y : BitVec 64) :
    archMinWord .fallback x y = .error "not implemented" := rfl

theorem archMaxWord_armInfinity (x : BitVec 64) :
    archMaxWord .arm64 x (BitVec.ofNat 64 0x7ff0000000000000) = .ok 0x7ff0000000000000 := by
  simp [archMaxWord, dimAssembly]

theorem archMinWord_armInfinity (x : BitVec 64) :
    archMinWord .arm64 x (BitVec.ofNat 64 0xfff0000000000000) = .ok 0xfff0000000000000 := by
  simp [archMinWord, dimAssembly]

theorem archMaxWord_armNaN :
    archMaxWord .arm64 (BitVec.ofNat 64 0x7ff8000000000001) (BitVec.ofNat 64 0xfff0000000000042)
      = .ok 0xfff8000000000042 := rfl

theorem archMinWord_armNaN :
    archMinWord .arm64 (BitVec.ofNat 64 0x7ff0000000000042) (BitVec.ofNat 64 0xfff8000000000001)
      = .ok 0x7ff8000000000042 := rfl

end GoLean.Math
