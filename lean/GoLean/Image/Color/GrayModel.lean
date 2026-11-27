import Init.Data.Nat.Bitwise.Lemmas

namespace GoLean.Image.Color

inductive GrayInput where
  | gray (value : Nat)
  | other (callback : Option (Nat × Nat × Nat))
  | nilColor

/-- Gray inputs bypass RGBA. Other inputs are converted with uint32
    multiplication and addition, then the highest byte of the sum is returned.
    The result records the callback count and distinguishes return from panic.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/color.go#L244-L261 -/
def grayModel : GrayInput → Option Nat × Nat
  | .gray value => (some value, 0)
  | .other callback =>
    (callback.map fun (r, g, b) =>
      ((19595 * r + 38470 * g + 7471 * b + (1 <<< 15)) % 2 ^ 32) >>> 24, 1)
  | .nilColor => (none, 0)

theorem grayModel_gray (value : Nat) :
    grayModel (.gray value) = (some value, 0) := rfl

theorem grayModel_return (r g b : Nat) :
    grayModel (.other (some (r, g, b))) =
      (some ((19595*r + 38470*g + 7471*b + 32768) / 16777216 % 256), 1) := by
  simp only [grayModel, Option.map_some, Nat.shiftRight_eq_div_pow, Nat.one_shiftLeft]
  change (some (((19595*r + 38470*g + 7471*b + 32768) % (16777216 * 256)) / 16777216), 1) = _
  rw [Nat.mod_mul_right_div_self]

theorem grayModel_bounded (r g b : Nat) (hr : r ≤ 65535) (hg : g ≤ 65535) (hb : b ≤ 65535) :
    grayModel (.other (some (r, g, b))) =
      (some ((19595*r + 38470*g + 7471*b + 32768) / 16777216), 1) := by
  rw [grayModel_return, Nat.mod_eq_of_lt (by omega)]

theorem grayModel_panic : grayModel (.other none) = (none, 1) := rfl

theorem grayModel_nil : grayModel .nilColor = (none, 0) := rfl

end GoLean.Image.Color
