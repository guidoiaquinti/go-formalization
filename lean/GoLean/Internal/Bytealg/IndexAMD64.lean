import Init.Data.BitVec.Basic

namespace GoLean.Internal.Bytealg

/-- Go 1.27.0 src/internal/bytealg/index_amd64.go:23–26. The addition
    wraps at native word width before signed division truncates toward zero. -/
def amd64Cutover (processed : BitVec 64) : Int :=
  ((processed + 16).toInt).tdiv 8

/-- Go 1.27.0 src/internal/bytealg/index_amd64.go:11–17.
    The initialized MaxLen overwrites that variable's old value. -/
def amd64IndexInit (hasAVX2 : Bool) : Nat :=
  if hasAVX2 then 63 else 31

theorem amd64Cutover_initial : amd64Cutover 0 = 2 := rfl

theorem amd64Cutover_truncation :
    amd64Cutover (BitVec.ofInt 64 (-17)) = 0 := rfl

theorem amd64Cutover_overflow :
    amd64Cutover 0x7fffffffffffffff = -1152921504606846974 := rfl

theorem amd64IndexInit_avx2 : amd64IndexInit true = 63 := rfl

theorem amd64IndexInit_base : amd64IndexInit false = 31 := rfl

end GoLean.Internal.Bytealg
