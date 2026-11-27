import GoLean.Encoding.Binary.EndianPut

namespace GoLean.Encoding.Binary

structure AppendWordResult where
  output : List (BitVec 8)
  originalBacking : List (BitVec 8)
  reused : Bool

def appendWord (bigEndian : Bool) (width word length : Nat)
    (backing : List (BitVec 8)) : AppendWordResult :=
  let result := (putWord bigEndian width word (List.replicate width 0)).2
  let combined := backing.take length ++ result
  let reused := length + width ≤ backing.length
  ⟨combined, if reused then combined ++ backing.drop (length + width) else backing, reused⟩

def bigAppendUint16 (word : BitVec 16) := appendWord true 2 word.toNat
def bigAppendUint32 (word : BitVec 32) := appendWord true 4 word.toNat
def bigAppendUint64 (word : BitVec 64) := appendWord true 8 word.toNat
def littleAppendUint16 (word : BitVec 16) := appendWord false 2 word.toNat
def littleAppendUint32 (word : BitVec 32) := appendWord false 4 word.toNat
def littleAppendUint64 (word : BitVec 64) := appendWord false 8 word.toNat

theorem appendWord_detached (bigEndian : Bool) (width word length : Nat)
    (backing : List (BitVec 8)) (short : backing.length < length + width) :
    (appendWord bigEndian width word length backing).originalBacking = backing := by
  simp [appendWord, Nat.not_le.mpr short]

theorem appendWord_reused (bigEndian : Bool) (width word length : Nat)
    (backing : List (BitVec 8)) :
    (appendWord bigEndian width word length backing).reused = decide (length + width ≤ backing.length) :=
  rfl

theorem bigAppendUint16_reused (word : BitVec 16) (length : Nat)
    (backing : List (BitVec 8)) :
    (bigAppendUint16 word length backing).reused = decide (length + 2 ≤ backing.length) :=
  appendWord_reused true 2 word.toNat length backing

theorem bigAppendUint16_detached (word : BitVec 16) (length : Nat)
    (backing : List (BitVec 8)) (short : backing.length < length + 2) :
    (bigAppendUint16 word length backing).originalBacking = backing :=
  appendWord_detached true 2 word.toNat length backing short

theorem bigAppendUint32_reused (word : BitVec 32) (length : Nat)
    (backing : List (BitVec 8)) :
    (bigAppendUint32 word length backing).reused = decide (length + 4 ≤ backing.length) :=
  appendWord_reused true 4 word.toNat length backing

theorem bigAppendUint32_detached (word : BitVec 32) (length : Nat)
    (backing : List (BitVec 8)) (short : backing.length < length + 4) :
    (bigAppendUint32 word length backing).originalBacking = backing :=
  appendWord_detached true 4 word.toNat length backing short

theorem bigAppendUint64_reused (word : BitVec 64) (length : Nat)
    (backing : List (BitVec 8)) :
    (bigAppendUint64 word length backing).reused = decide (length + 8 ≤ backing.length) :=
  appendWord_reused true 8 word.toNat length backing

theorem bigAppendUint64_detached (word : BitVec 64) (length : Nat)
    (backing : List (BitVec 8)) (short : backing.length < length + 8) :
    (bigAppendUint64 word length backing).originalBacking = backing :=
  appendWord_detached true 8 word.toNat length backing short

theorem littleAppendUint16_reused (word : BitVec 16) (length : Nat)
    (backing : List (BitVec 8)) :
    (littleAppendUint16 word length backing).reused = decide (length + 2 ≤ backing.length) :=
  appendWord_reused false 2 word.toNat length backing

theorem littleAppendUint16_detached (word : BitVec 16) (length : Nat)
    (backing : List (BitVec 8)) (short : backing.length < length + 2) :
    (littleAppendUint16 word length backing).originalBacking = backing :=
  appendWord_detached false 2 word.toNat length backing short

theorem littleAppendUint32_reused (word : BitVec 32) (length : Nat)
    (backing : List (BitVec 8)) :
    (littleAppendUint32 word length backing).reused = decide (length + 4 ≤ backing.length) :=
  appendWord_reused false 4 word.toNat length backing

theorem littleAppendUint32_detached (word : BitVec 32) (length : Nat)
    (backing : List (BitVec 8)) (short : backing.length < length + 4) :
    (littleAppendUint32 word length backing).originalBacking = backing :=
  appendWord_detached false 4 word.toNat length backing short

theorem littleAppendUint64_reused (word : BitVec 64) (length : Nat)
    (backing : List (BitVec 8)) :
    (littleAppendUint64 word length backing).reused = decide (length + 8 ≤ backing.length) :=
  appendWord_reused false 8 word.toNat length backing

theorem littleAppendUint64_detached (word : BitVec 64) (length : Nat)
    (backing : List (BitVec 8)) (short : backing.length < length + 8) :
    (littleAppendUint64 word length backing).originalBacking = backing :=
  appendWord_detached false 8 word.toNat length backing short

end GoLean.Encoding.Binary
