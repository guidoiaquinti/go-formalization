import GoLean.Image.Color.RGBAModel

namespace GoLean.Image.Color

private def convert (narrow : Nat → Nat) (r g b a : Nat) : Channels :=
  if a = 65535 then
    (narrow r, narrow g, narrow b, narrow 65535)
  else if a = 0 then
    (0, 0, 0, 0)
  else
    let unpremultiply := fun x => ((x * 65535) % 4294967296) / a
    (narrow (unpremultiply r), narrow (unpremultiply g), narrow (unpremultiply b), narrow a)

/-- Convert a Color callback to NRGBA. A callback's uint32 products wrap
    before division. A same-type value bypasses the callback entirely.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/color.go#L192-L208 -/
def nrgbaModel : RGBAInput → Option Channels × Nat
  | .sameType value => (some value, 0)
  | .other callback => (callback.map fun (r, g, b, a) => convert (fun x => (x >>> 8) &&& 255) r g b a, 1)
  | .nilColor => (none, 0)

/-- Convert a Color callback to NRGBA64 using the same uint32 products
    and taking the low 16 bits of the resulting words.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/color.go#L210-L226 -/
def nrgba64Model : RGBAInput → Option Channels × Nat
  | .sameType value => (some value, 0)
  | .other callback => (callback.map fun (r, g, b, a) => convert (fun x => x &&& 65535) r g b a, 1)
  | .nilColor => (none, 0)

private def expected (narrow : Nat → Nat) (r g b a : Nat) : Channels :=
  if a = 0 then
    (0, 0, 0, 0)
  else
    let scale := fun x => if a = 65535 then x else x * 65535 % 4294967296 / a
    (narrow (scale r), narrow (scale g), narrow (scale b), narrow a)

theorem nrgbaModel_same (c : Channels) :
    nrgbaModel (.sameType c) = (some c, 0) := rfl

theorem nrgbaModel_return (r g b a : Nat) :
    nrgbaModel (.other (some (r, g, b, a))) =
      (some (expected (fun x => x / 256 % 256) r g b a), 1) := by
  change (some (convert (fun x => (x >>> 8) &&& (2 ^ 8 - 1)) r g b a), 1) = _
  simp only [Nat.and_two_pow_sub_one_eq_mod, Nat.shiftRight_eq_div_pow]
  by_cases zero : a = 0
  · simp [zero, convert, expected]
  by_cases fullAlpha : a = 65535
  · simp [fullAlpha, convert, expected]
  · simp [zero, fullAlpha, convert, expected]

theorem nrgbaModel_panic :
    nrgbaModel (.other none) = (none, 1) ∧ nrgbaModel .nilColor = (none, 0) :=
  ⟨rfl, rfl⟩

theorem nrgba64Model_same (c : Channels) :
    nrgba64Model (.sameType c) = (some c, 0) := rfl

theorem nrgba64Model_return (r g b a : Nat) :
    nrgba64Model (.other (some (r, g, b, a))) =
      (some (expected (fun x => x % 65536) r g b a), 1) := by
  change (some (convert (fun x => x &&& (2 ^ 16 - 1)) r g b a), 1) = _
  simp only [Nat.and_two_pow_sub_one_eq_mod]
  by_cases zero : a = 0
  · simp [zero, convert, expected]
  by_cases fullAlpha : a = 65535
  · simp [fullAlpha, convert, expected]
  · simp [zero, fullAlpha, convert, expected]

theorem nrgba64Model_panic :
    nrgba64Model (.other none) = (none, 1) ∧ nrgba64Model .nilColor = (none, 0) :=
  ⟨rfl, rfl⟩

end GoLean.Image.Color
