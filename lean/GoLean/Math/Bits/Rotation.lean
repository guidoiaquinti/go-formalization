import Init.Data.BitVec.Lemmas

namespace GoLean.Math.Bits

/-- The signed Go int count may include MinInt. Euclidean remainder agrees with
    masking uint(k) by w-1 at the supported power-of-two widths, even for k < 0.
    A zero count shifts the second term by the whole width. -/
private def rotateWord (x : BitVec w) (k : Int) : BitVec w :=
  let s := (k % (w : Int)).toNat
  (x <<< s) ||| (x >>> (w - s))

/-- The low s result bits come from the top of x. The other bits move left by s. -/
def IsLeftRotation (x y : BitVec w) (k : Int) : Prop :=
  let s := (k % (w : Int)).toNat
  ∀ i < w, y.getLsbD i = x.getLsbD (if i < s then w - s + i else i - s)

private theorem rotate_spec (x : BitVec w) (k : Int) (hw : 0 < w) :
    IsLeftRotation x (rotateWord x k) k := by
  let s := (k % (w : Int)).toNat
  have hpos : (0 : Int) < w := by omega
  have hs : s < w := by
    have := Int.emod_lt_of_pos k hpos
    omega
  have heq : rotateWord x k = x.rotateLeft s := by
    simp [rotateWord, BitVec.rotateLeft, BitVec.rotateLeftAux, s, Nat.mod_eq_of_lt hs]
  intro i hi
  rw [heq, BitVec.getLsbD_rotateLeft_of_le hs]
  by_cases h : i < s <;> simp [h, hi, s] at *

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L187 -/
def rotateLeft8 (x : BitVec 8) (k : Int) : BitVec 8 := rotateWord x k

theorem rotateLeft8_spec (x : BitVec 8) (k : Int) : IsLeftRotation x (rotateLeft8 x k) k :=
  rotate_spec x k (by decide)

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L197 -/
def rotateLeft16 (x : BitVec 16) (k : Int) : BitVec 16 := rotateWord x k

theorem rotateLeft16_spec (x : BitVec 16) (k : Int) : IsLeftRotation x (rotateLeft16 x k) k :=
  rotate_spec x k (by decide)

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L207 -/
def rotateLeft32 (x : BitVec 32) (k : Int) : BitVec 32 := rotateWord x k

theorem rotateLeft32_spec (x : BitVec 32) (k : Int) : IsLeftRotation x (rotateLeft32 x k) k :=
  rotate_spec x k (by decide)

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L217 -/
def rotateLeft64 (x : BitVec 64) (k : Int) : BitVec 64 := rotateWord x k

theorem rotateLeft64_spec (x : BitVec 64) (k : Int) : IsLeftRotation x (rotateLeft64 x k) k :=
  rotate_spec x k (by decide)

/-- RotateLeft on the pinned 64-bit Go profile.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L176 -/
def rotateLeft (x : BitVec 64) (k : Int) : BitVec 64 := rotateLeft64 x k

theorem rotateLeft_spec (x : BitVec 64) (k : Int) : IsLeftRotation x (rotateLeft x k) k :=
  rotate_spec x k (by decide)

end GoLean.Math.Bits
