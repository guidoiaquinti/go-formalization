import Init.Data.BitVec.Lemmas
import Init.Omega

namespace GoLean.Compress.Flate

abbrev HCode := BitVec 32

def newHCode (code : BitVec 16) (length : BitVec 8) : HCode :=
  length.zeroExtend 32 ||| (code.zeroExtend 32 <<< 8)

def hCodeLen (word : HCode) : BitVec 8 := word.truncate 8

def hCodeCode64 (word : HCode) : BitVec 64 := (word >>> 8).zeroExtend 64

def hCodeZero (word : HCode) : Bool := word == 0

def hCodeSet (_previous : HCode) (code : BitVec 16) (length : BitVec 8) : HCode :=
  newHCode code length

private theorem bitPositions {p : Nat → Prop} (head : p 0)
    (more : ∀ i < n, p (i + 1)) : ∀ i < n + 1, p i := by
  intro i hi
  cases i with
  | zero => exact head
  | succ i => exact more i (by omega)

theorem newHCode_parts (code : BitVec 16) (length : BitVec 8) :
    hCodeLen (newHCode code length) = length ∧
      hCodeCode64 (newHCode code length) = code.zeroExtend 64 := by
  constructor
  · rw [BitVec.eq_of_getLsbD_eq_iff]
    iterate 8
      apply bitPositions
      · simp +decide [hCodeLen, newHCode]
    omega
  · rw [BitVec.eq_of_getLsbD_eq_iff]
    iterate 64
      apply bitPositions
      · simp +decide [hCodeCode64, newHCode]
    omega

theorem hCodeLen_bits (word : HCode) :
    (hCodeLen word).zeroExtend 32 = word &&& 255 := by
  rw [BitVec.eq_of_getLsbD_eq_iff]
  iterate 32
    apply bitPositions
    · simp +decide [hCodeLen]
  omega

theorem hCodeCode64_bits (word : HCode) :
    (hCodeCode64 word <<< 8).truncate 32 = word &&& 0xffffff00 := by
  rw [BitVec.eq_of_getLsbD_eq_iff]
  iterate 32
    apply bitPositions
    · simp +decide [hCodeCode64]
  omega

theorem hCodeZero_iff (word : HCode) : hCodeZero word = true ↔ word = 0 := by
  simp [hCodeZero]

theorem hCodeSet_overwrites (old : HCode) (code : BitVec 16) (length : BitVec 8) :
    hCodeLen (hCodeSet old code length) = length ∧
      hCodeCode64 (hCodeSet old code length) = code.zeroExtend 64 :=
  newHCode_parts code length

end GoLean.Compress.Flate
