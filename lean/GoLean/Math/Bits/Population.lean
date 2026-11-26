import GoLean.Math.Bits.OnesCount8
import Init.Data.BitVec.Lemmas

namespace GoLean.Math.Bits

/-- Summing the two half-word counts models the disjoint byte-table sums at
    16 and 32 bits and the adjacent-field sums of the 64-bit parallel algorithm.
    The intermediates are nonnegative and their sum is at most the word width. -/
private def sumHalves (count : BitVec w → Nat) (x : BitVec (w + w)) : Nat :=
  count (x.extractLsb' w w) + count (x.extractLsb' 0 w)

private theorem sumHalves_spec (count : BitVec w → Nat)
    (h : ∀ b, count b = b.cpop.toNat) (x : BitVec (w + w)) :
    sumHalves count x = x.cpop.toNat := by
  unfold sumHalves
  rw [h, h, ← BitVec.toNat_cpop_append, BitVec.extractLsb'_append_extractLsb']

private theorem byte_spec : ∀ x : BitVec 8, onesCount8 x = x.cpop.toNat := by decide

/-- Count set positions independently of the model's byte divisions. -/
private theorem cpop_count (x : BitVec w) :
    x.cpop.toNat = (List.range w).countP x.getLsbD := by
  rw [BitVec.toNat_cpop]
  suffices h : ∀ n, x.cpopNatRec n 0 = (List.range n).countP x.getLsbD from h w
  intro n
  induction n with
  | zero => rfl
  | succ n ih =>
    rw [BitVec.cpopNatRec_succ, BitVec.cpopNatRec_eq, List.range_succ, List.countP_append, ih]
    cases h : x.getLsbD n <;> simp [h]

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L130 -/
def onesCount16 (x : BitVec 16) : Nat := sumHalves onesCount8 x

private theorem pop16_eq (x : BitVec 16) : onesCount16 x = x.cpop.toNat :=
  sumHalves_spec _ byte_spec x

theorem onesCount16_spec (x : BitVec 16) :
    onesCount16 x = (List.range 16).countP x.getLsbD :=
  (pop16_eq x).trans (cpop_count x)

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L135 -/
def onesCount32 (x : BitVec 32) : Nat := sumHalves onesCount16 x

private theorem pop32_eq (x : BitVec 32) : onesCount32 x = x.cpop.toNat :=
  sumHalves_spec _ pop16_eq x

theorem onesCount32_spec (x : BitVec 32) :
    onesCount32 x = (List.range 32).countP x.getLsbD :=
  (pop32_eq x).trans (cpop_count x)

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L140 -/
def onesCount64 (x : BitVec 64) : Nat := sumHalves onesCount32 x

private theorem pop64_eq (x : BitVec 64) : onesCount64 x = x.cpop.toNat :=
  sumHalves_spec _ pop32_eq x

theorem onesCount64_spec (x : BitVec 64) :
    onesCount64 x = (List.range 64).countP x.getLsbD :=
  (pop64_eq x).trans (cpop_count x)

/-- OnesCount on the pinned 64-bit Go profile.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L117 -/
def onesCount := onesCount64

theorem onesCount_spec (x : BitVec 64) :
    onesCount x = (List.range 64).countP x.getLsbD := onesCount64_spec x

end GoLean.Math.Bits
