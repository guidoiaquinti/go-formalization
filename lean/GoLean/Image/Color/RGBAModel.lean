import Init.Data.Nat.Bitwise.Lemmas

namespace GoLean.Image.Color

abbrev Channels := Nat × Nat × Nat × Nat

private def mapChannels (f : Nat → Nat) (c : Channels) : Channels :=
  (f c.1, f c.2.1, f c.2.2.1, f c.2.2.2)

inductive RGBAInput where
  | sameType (value : Channels)
  | other (callback : Option Channels)
  | nilColor

/-- Go preserves RGBA values. Otherwise it shifts each callback word by eight
    and takes the low byte. None represents panic.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/color.go#L176-L182 -/
def rgbaModel : RGBAInput → Option Channels × Nat
  | .sameType value => (some value, 0)
  | .other callback => (callback.map (mapChannels fun x => (x >>> 8) &&& 255), 1)
  | .nilColor => (none, 0)

/-- Go preserves RGBA64 values. Otherwise it takes the low 16 bits of each
    callback word.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/color.go#L184-L190 -/
def rgba64Model : RGBAInput → Option Channels × Nat
  | .sameType value => (some value, 0)
  | .other callback => (callback.map (mapChannels fun x => x &&& 65535), 1)
  | .nilColor => (none, 0)

theorem rgbaModel_same (c : Channels) :
    rgbaModel (.sameType c) = (some c, 0) := rfl

theorem rgbaModel_return (r g b a : Nat) :
    rgbaModel (.other (some (r, g, b, a))) =
      (some (r / 256 % 256, g / 256 % 256, b / 256 % 256, a / 256 % 256), 1) := by
  change (some ((r >>> 8) &&& (2 ^ 8 - 1), (g >>> 8) &&& (2 ^ 8 - 1),
    (b >>> 8) &&& (2 ^ 8 - 1), (a >>> 8) &&& (2 ^ 8 - 1)), 1) = _
  simp only [Nat.shiftRight_eq_div_pow, Nat.and_two_pow_sub_one_eq_mod]

theorem rgbaModel_panic :
    rgbaModel (.other none) = (none, 1) ∧ rgbaModel .nilColor = (none, 0) := by
  exact ⟨rfl, rfl⟩

theorem rgba64Model_same (c : Channels) :
    rgba64Model (.sameType c) = (some c, 0) := rfl

theorem rgba64Model_return (r g b a : Nat) :
    rgba64Model (.other (some (r, g, b, a))) =
      (some (r % 65536, g % 65536, b % 65536, a % 65536), 1) := by
  change (some (r &&& (2 ^ 16 - 1), g &&& (2 ^ 16 - 1),
    b &&& (2 ^ 16 - 1), a &&& (2 ^ 16 - 1)), 1) = _
  simp only [Nat.and_two_pow_sub_one_eq_mod]

theorem rgba64Model_panic :
    rgba64Model (.other none) = (none, 1) ∧ rgba64Model .nilColor = (none, 0) := by
  exact ⟨rfl, rfl⟩

end GoLean.Image.Color
