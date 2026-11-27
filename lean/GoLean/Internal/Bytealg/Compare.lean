import GoLean.Bytes.Compare

namespace GoLean.Internal.Bytealg

/-- Go 1.27.0 src/internal/bytealg/compare_native.go:12, with byte semantics
    documented by the generic implementation in compare_generic.go. -/
def compare := GoLean.Bytes.compare

/-- Go 1.27.0 src/internal/bytealg/compare_native.go:23. -/
def runtimeCmpstring (a b : List (BitVec 8)) : Int :=
  compare (some (a, a.length)) (some (b, b.length))

/-- Go 1.27.0 src/internal/bytealg/compare_native.go:14–16. -/
def compareString := runtimeCmpstring

/-- Go 1.27.0 src/internal/bytealg/equal_generic.go:14–18. -/
def equal := GoLean.Bytes.equal

theorem compare_spec (a b : GoLean.Bytes.Slice) :
    let x := (GoLean.Bytes.contents a).map fun byte => (byte.toNat : Int)
    let y := (GoLean.Bytes.contents b).map fun byte => (byte.toNat : Int)
    (compare a b = -1 ↔ x < y) ∧
    (compare a b = 0 ↔ x = y) ∧
    (compare a b = 1 ↔ y < x) :=
  GoLean.Bytes.compare_spec a b

theorem runtimeCmpstring_spec (a b : List (BitVec 8)) :
    let x := a.map fun byte => (byte.toNat : Int)
    let y := b.map fun byte => (byte.toNat : Int)
    (runtimeCmpstring a b = -1 ↔ x < y) ∧
    (runtimeCmpstring a b = 0 ↔ x = y) ∧
    (runtimeCmpstring a b = 1 ↔ y < x) :=
  compare_spec (some (a, a.length)) (some (b, b.length))

theorem compareString_spec (a b : List (BitVec 8)) :
    let x := a.map fun byte => (byte.toNat : Int)
    let y := b.map fun byte => (byte.toNat : Int)
    (compareString a b = -1 ↔ x < y) ∧
    (compareString a b = 0 ↔ x = y) ∧
    (compareString a b = 1 ↔ y < x) :=
  runtimeCmpstring_spec a b

theorem equal_spec (a b : GoLean.Bytes.Slice) :
    equal a b = true ↔ GoLean.Bytes.contents a = GoLean.Bytes.contents b :=
  GoLean.Bytes.equal_spec a b

end GoLean.Internal.Bytealg
