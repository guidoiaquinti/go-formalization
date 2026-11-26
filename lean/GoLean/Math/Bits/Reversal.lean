import Init.Data.BitVec.Lemmas
import Init.Data.BitVec.Decidable

namespace GoLean.Math.Bits

/-- Each output bit keeps its offset within its byte and moves to the opposite byte. -/
def IsByteReverse (x y : BitVec w) : Prop :=
  ∀ i < w, y.getLsbD i = x.getLsbD ((w / 8 - 1 - i / 8) * 8 + i % 8)

def IsBitReverse (x y : BitVec w) : Prop :=
  ∀ i < w, y.getLsbD i = x.getLsbD (w - 1 - i)

/-- Split a finite set of positions while leaving the input word symbolic. -/
private theorem positions_succ {p : Nat → Prop} (first : p 0)
    (rest : ∀ i < n, p (i + 1)) : ∀ i < n + 1, p i := by
  intro i hi
  cases i with
  | zero => exact first
  | succ i => exact rest i (by omega)

/-- Exchange adjacent groups of s bits selected by Go's alternating mask. -/
private def exchange (x mask : BitVec w) (s : Nat) : BitVec w :=
  ((x >>> s) &&& mask) ||| ((x &&& mask) <<< s)

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L283 -/
def reverseBytes32 (x : BitVec 32) : BitVec 32 :=
  let x := exchange x 0x00ff00ff 8
  (x >>> 16) ||| (x <<< 16)

theorem reverseBytes32_spec (x : BitVec 32) : IsByteReverse x (reverseBytes32 x) := by
  unfold IsByteReverse
  iterate 32
    apply positions_succ
    · simp +decide [reverseBytes32, exchange]
  omega

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L292 -/
def reverseBytes64 (x : BitVec 64) : BitVec 64 :=
  let x := exchange x 0x00ff00ff00ff00ff 8
  let x := exchange x 0x0000ffff0000ffff 16
  (x >>> 32) ||| (x <<< 32)

theorem reverseBytes64_spec (x : BitVec 64) : IsByteReverse x (reverseBytes64 x) := by
  unfold IsByteReverse
  iterate 64
    apply positions_succ
    · simp +decide [reverseBytes64, exchange]
  omega

/-- ReverseBytes on the pinned 64-bit Go profile.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L266 -/
def reverseBytes := reverseBytes64

theorem reverseBytes_spec (x : BitVec 64) : IsByteReverse x (reverseBytes x) := reverseBytes64_spec x

/-- One iteration of the pinned rev8tab generator moves the low bit to the accumulator.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/make_tables.go#L77-L83 -/
private def reverseInto : Nat → BitVec 8 → BitVec 8 → BitVec 8
  | 0, _, acc => acc
  | n + 1, x, acc => reverseInto n (x >>> 1) ((acc <<< 1) ||| (x &&& 1))

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L234 -/
def reverse8 (x : BitVec 8) : BitVec 8 := reverseInto 8 x 0

private theorem reverse8_eq : ∀ x : BitVec 8, reverse8 x = x.reverse := by decide

theorem reverse8_spec (x : BitVec 8) : IsBitReverse x (reverse8 x) := by
  intro i hi
  rw [reverse8_eq, BitVec.getLsbD_reverse, BitVec.getMsbD_eq_getLsbD]
  simp [hi]

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L239 -/
def reverse16 (x : BitVec 16) : BitVec 16 :=
  reverse8 (x.extractLsb' 0 8) ++ reverse8 (x.extractLsb' 8 8)

theorem reverse16_spec (x : BitVec 16) : IsBitReverse x (reverse16 x) := by
  have halves := BitVec.reverse_append (x := x.extractLsb' 8 8) (y := x.extractLsb' 0 8)
  rw [BitVec.extractLsb'_append_extractLsb', BitVec.cast_eq] at halves
  intro i hi
  rw [reverse16, reverse8_eq, reverse8_eq, ← show x.reverse = _ from halves,
    BitVec.getLsbD_reverse, BitVec.getMsbD_eq_getLsbD]
  simp [hi]

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L244 -/
def reverse32 (x : BitVec 32) : BitVec 32 :=
  let x := exchange x 0x55555555 1
  let x := exchange x 0x33333333 2
  let x := exchange x 0x0f0f0f0f 4
  reverseBytes32 x

theorem reverse32_spec (x : BitVec 32) : IsBitReverse x (reverse32 x) := by
  unfold IsBitReverse
  iterate 32
    apply positions_succ
    · simp +decide [reverse32, reverseBytes32, exchange]
  omega

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L253 -/
def reverse64 (x : BitVec 64) : BitVec 64 :=
  let x := exchange x 0x5555555555555555 1
  let x := exchange x 0x3333333333333333 2
  let x := exchange x 0x0f0f0f0f0f0f0f0f 4
  reverseBytes64 x

theorem reverse64_spec (x : BitVec 64) : IsBitReverse x (reverse64 x) := by
  unfold IsBitReverse
  iterate 64
    apply positions_succ
    · simp +decide [reverse64, reverseBytes64, exchange]
  omega

/-- Reverse on the pinned 64-bit Go profile.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L226 -/
def reverse := reverse64

theorem reverse_spec (x : BitVec 64) : IsBitReverse x (reverse x) := reverse64_spec x

end GoLean.Math.Bits
