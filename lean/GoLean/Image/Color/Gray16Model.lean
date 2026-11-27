import GoLean.Image.Color.GrayModel

namespace GoLean.Image.Color

/-- Gray16 inputs bypass RGBA. Other inputs are converted with uint32
    multiplication and addition, then the high 16 bits of the sum is returned.
    The result records the callback count and distinguishes return from panic.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/color.go#L263-L277 -/
def gray16Model : GrayInput → Option Nat × Nat
  | .gray value => (some value, 0)
  | .other callback =>
    (callback.map fun (r, g, b) =>
      ((19595 * r + 38470 * g + 7471 * b + (1 <<< 15)) % 2 ^ 32) >>> 16, 1)
  | .nilColor => (none, 0)

theorem gray16Model_gray (value : Nat) :
    gray16Model (.gray value) = (some value, 0) := rfl

theorem gray16Model_return (r g b : Nat) :
    gray16Model (.other (some (r, g, b))) =
      (some ((19595*r + 38470*g + 7471*b + 32768) / 65536 % 65536), 1) := by
  simp only [gray16Model, Option.map_some, Nat.shiftRight_eq_div_pow, Nat.one_shiftLeft]
  change (some (((19595*r + 38470*g + 7471*b + 32768) % (65536 * 65536)) / 65536), 1) = _
  rw [Nat.mod_mul_right_div_self]

theorem gray16Model_bounded (r g b : Nat) (hr : r ≤ 65535) (hg : g ≤ 65535) (hb : b ≤ 65535) :
    gray16Model (.other (some (r, g, b))) =
      (some ((19595*r + 38470*g + 7471*b + 32768) / 65536), 1) := by
  rw [gray16Model_return, Nat.mod_eq_of_lt (by omega)]

theorem gray16Model_panic : gray16Model (.other none) = (none, 1) := rfl

theorem gray16Model_nil : gray16Model .nilColor = (none, 0) := rfl

end GoLean.Image.Color
