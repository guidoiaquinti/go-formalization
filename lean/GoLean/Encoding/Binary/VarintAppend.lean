import GoLean.Encoding.Binary.EndianAppend

namespace GoLean.Encoding.Binary

def uvarintBytes (word : Nat) : List (BitVec 8) :=
  if word < 128 then [BitVec.ofNat 8 word]
  else BitVec.ofNat 8 (128 + word % 128) :: uvarintBytes (word / 128)
termination_by word
decreasing_by
  exact Nat.div_lt_self (by omega) (by decide)

def signedVarintWord (x : BitVec 64) : Nat :=
  let shifted := x <<< 1
  if x.toInt < 0 then (~~~shifted).toNat else shifted.toNat

private def appendBytes (payload : List (BitVec 8)) (length : Nat)
    (backing : List (BitVec 8)) : AppendWordResult :=
  let original := backing.take length
  let written := payload.take (backing.length - length)
  ⟨original ++ payload, original ++ written ++ backing.drop (length + written.length),
    decide (length + payload.length ≤ backing.length)⟩

def appendUvarint (x : BitVec 64) (length : Nat) (backing : List (BitVec 8)) :
    AppendWordResult :=
  appendBytes (uvarintBytes x.toNat) length backing

def appendVarint (x : BitVec 64) (length : Nat) (backing : List (BitVec 8)) :
    AppendWordResult :=
  appendBytes (uvarintBytes (signedVarintWord x)) length backing

theorem appendUvarint_payload (x : BitVec 64) (length : Nat) (backing : List (BitVec 8)) :
    (appendUvarint x length backing).output = backing.take length ++ uvarintBytes x.toNat :=
  rfl

theorem appendVarint_payload (x : BitVec 64) (length : Nat) (backing : List (BitVec 8)) :
    (appendVarint x length backing).output =
      backing.take length ++ uvarintBytes (signedVarintWord x) :=
  rfl

theorem appendUvarint_reused (x : BitVec 64) (length : Nat) (backing : List (BitVec 8)) :
    (appendUvarint x length backing).reused =
      decide (length + (uvarintBytes x.toNat).length ≤ backing.length) :=
  rfl

theorem appendVarint_reused (x : BitVec 64) (length : Nat) (backing : List (BitVec 8)) :
    (appendVarint x length backing).reused =
      decide (length + (uvarintBytes (signedVarintWord x)).length ≤ backing.length) :=
  rfl

theorem appendUvarint_backing (x : BitVec 64) (length : Nat) (backing : List (BitVec 8)) :
    (appendUvarint x length backing).originalBacking =
      let written := (uvarintBytes x.toNat).take (backing.length - length)
      backing.take length ++ written ++ backing.drop (length + written.length) :=
  rfl

theorem appendVarint_backing (x : BitVec 64) (length : Nat) (backing : List (BitVec 8)) :
    (appendVarint x length backing).originalBacking =
      let written := (uvarintBytes (signedVarintWord x)).take (backing.length - length)
      backing.take length ++ written ++ backing.drop (length + written.length) :=
  rfl

end GoLean.Encoding.Binary
