import Init.Data.BitVec.Basic

namespace GoLean.Internal.Bytealg

/-- Go 1.27.0 src/internal/bytealg/count_native.go:16–23. -/
def countGeneric (bytes : List (BitVec 8)) (target : BitVec 8) : Nat :=
  bytes.foldl (fun count byte => if byte == target then count + 1 else count) 0

/-- Go counts raw string bytes, using the same equality as the slice loop. -/
def countGenericString (bytes : List (BitVec 8)) (target : BitVec 8) : Nat :=
  countGeneric bytes target

/-- Byte-count result of the native slice entry point. -/
def count (bytes : List (BitVec 8)) (target : BitVec 8) : Nat :=
  countGeneric bytes target

/-- Byte-count result of the native string entry point. -/
def countString (bytes : List (BitVec 8)) (target : BitVec 8) : Nat :=
  countGenericString bytes target

private theorem countFold_spec (bytes : List (BitVec 8)) (target : BitVec 8) (initial : Nat) :
    bytes.foldl (fun count byte => if byte == target then count + 1 else count) initial =
      initial + (bytes.filter (· == target)).length := by
  induction bytes generalizing initial with
  | nil => simp
  | cons byte bytes ih =>
    simp only [List.foldl_cons, List.filter_cons]
    rw [ih]
    split <;> simp_all <;> omega

theorem countGeneric_spec (bytes : List (BitVec 8)) (target : BitVec 8) :
    countGeneric bytes target = (bytes.filter (· == target)).length := by
  simpa [countGeneric] using countFold_spec bytes target 0

theorem countGenericString_spec (bytes : List (BitVec 8)) (target : BitVec 8) :
    countGenericString bytes target = (bytes.filter (· == target)).length :=
  countGeneric_spec bytes target

theorem count_spec (bytes : List (BitVec 8)) (target : BitVec 8) :
    count bytes target = (bytes.filter (· == target)).length :=
  countGeneric_spec bytes target

theorem countString_spec (bytes : List (BitVec 8)) (target : BitVec 8) :
    countString bytes target = (bytes.filter (· == target)).length :=
  countGenericString_spec bytes target

end GoLean.Internal.Bytealg
