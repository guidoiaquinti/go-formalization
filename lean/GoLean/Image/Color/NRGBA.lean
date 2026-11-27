import Init.Data.BitVec.Lemmas

namespace GoLean.Image.Color

private def widenByte (x : BitVec 8) : Nat :=
  let wide := x.zeroExtend 32
  (wide ||| (wide <<< 8)).toNat

private def premultiply (channel alpha : BitVec 8) : Nat :=
  (widenByte channel * alpha.toNat) / 255

/-- Go widens each color byte by repeating it, then multiplies by the original
    alpha byte and divides by 255. The largest intermediate product is
    65535 * 255, which fits in uint32. The output alpha is the repeated byte.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/color.go#L55-L76 -/
def nrgbaRGBA (r g b a : BitVec 8) : Nat × Nat × Nat × Nat :=
  (premultiply r a, premultiply g a, premultiply b a, widenByte a)

private theorem widenByte_spec : ∀ x : BitVec 8, widenByte x = x.toNat * 257 := by
  decide

private theorem premultiply_spec (x a : BitVec 8) :
    premultiply x a = x.toNat * a.toNat + 2 * (x.toNat * a.toNat) / 255 := by
  unfold premultiply
  rw [widenByte_spec]
  let p := x.toNat * a.toNat
  rw [show x.toNat * 257 * a.toNat = 257 * p by
        dsimp only [p]
        ac_rfl]
  change ((2 + 255) * p) / 255 = p + 2 * p / 255
  rw [Nat.add_mul, Nat.add_mul_div_left _ _ (by decide)]
  exact Nat.add_comm _ _

theorem nrgbaRGBA_spec (r g b a : BitVec 8) :
    nrgbaRGBA r g b a =
      (r.toNat * a.toNat + 2 * (r.toNat * a.toNat) / 255,
       g.toNat * a.toNat + 2 * (g.toNat * a.toNat) / 255,
       b.toNat * a.toNat + 2 * (b.toNat * a.toNat) / 255,
       a.toNat * 257) := by
  simp only [nrgbaRGBA, premultiply_spec, widenByte_spec]

end GoLean.Image.Color
