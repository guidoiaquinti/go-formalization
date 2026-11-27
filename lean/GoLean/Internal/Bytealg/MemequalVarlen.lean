import GoLean.Internal.Bytealg.RuntimeMemory

namespace GoLean.Internal.Bytealg

/-- Go 1.27.0 internal/bytealg/equal_native.go:21, and the arm64
    runtime.memequal_varlen entry in equal_arm64.s. The compiler puts
    the comparison size in closure metadata, then the entry compares
    that many bytes by delegating to the fixed-size runtime helper. -/
def runtimeMemequalVarlen (closureSize : Nat) (left right : List (BitVec 8)) : Bool :=
  runtimeMemequal closureSize left right

theorem runtimeMemequalVarlen_spec (closureSize : Nat) (left right : List (BitVec 8)) :
    runtimeMemequalVarlen closureSize left right = true ↔
      left.take closureSize = right.take closureSize :=
  runtimeMemequal_spec closureSize left right

theorem runtimeMemequalVarlen_zero (left right : List (BitVec 8)) :
    runtimeMemequalVarlen 0 left right = true :=
  runtimeMemequal_zero left right

end GoLean.Internal.Bytealg
