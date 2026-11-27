import Init.Data.Nat.Bitwise.Lemmas

namespace GoLean.Image.Color

inductive Alpha16Input where
  | alpha16 (value : Nat)
  | other (callback : Option Nat)
  | nilColor

/-- The Alpha16 fast path preserves the value. The fallback narrows the
    callback alpha word to 16 bits. None represents a panic, and the count
    records RGBA invocations.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/color.go#L236-L242 -/
def alpha16Model : Alpha16Input → Option Nat × Nat
  | .alpha16 value => (some value, 0)
  | .other callback => (callback.map fun a => a &&& 65535, 1)
  | .nilColor => (none, 0)

theorem alpha16Model_alpha16 (a : Nat) :
    alpha16Model (.alpha16 a) = (some a, 0) := rfl

theorem alpha16Model_return (a : Nat) :
    alpha16Model (.other (some a)) = (some (a % 65536), 1) := by
  change (some (a &&& (2 ^ 16 - 1)), 1) = _
  rw [Nat.and_two_pow_sub_one_eq_mod]

theorem alpha16Model_panic :
    alpha16Model (.other none) = (none, 1) := rfl

theorem alpha16Model_nil :
    alpha16Model .nilColor = (none, 0) := rfl

end GoLean.Image.Color
