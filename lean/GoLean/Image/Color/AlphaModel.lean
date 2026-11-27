import Init.Data.Nat.Bitwise.Lemmas

namespace GoLean.Image.Color

inductive AlphaInput where
  | alpha (value : Nat)
  | other (callback : Option Nat)
  | nilColor

/-- A missing result represents a panic. The count tracks RGBA invocations.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/color.go#L228-L234 -/
def alphaModel : AlphaInput → Option Nat × Nat
  | .alpha value => (some value, 0)
  | .other callback => (callback.map fun a => a >>> 8 &&& 255, 1)
  | .nilColor => (none, 0)

theorem alphaModel_alpha (a : Nat) :
    alphaModel (.alpha a) = (some a, 0) := rfl

theorem alphaModel_return (a : Nat) :
    alphaModel (.other (some a)) = (some (a / 256 % 256), 1) := by
  change (some (a >>> 8 &&& (2 ^ 8 - 1)), 1) = _
  rw [Nat.and_two_pow_sub_one_eq_mod, Nat.shiftRight_eq_div_pow]

theorem alphaModel_panic :
    alphaModel (.other none) = (none, 1) := rfl

theorem alphaModel_nil :
    alphaModel .nilColor = (none, 0) := rfl

end GoLean.Image.Color
