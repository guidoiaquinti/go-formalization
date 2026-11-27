import GoLean.Encoding.Binary.VarintAppend

namespace GoLean.Encoding.Binary

private def putPayload (payload backing : List (BitVec 8)) (length : Nat) :
    Option Nat × List (BitVec 8) :=
  (if payload.length ≤ length then some payload.length else none,
    payload.take length ++ backing.drop (min length payload.length))

def putUvarint (x : BitVec 64) (length : Nat) (backing : List (BitVec 8)) :
    Option Nat × List (BitVec 8) :=
  putPayload (uvarintBytes x.toNat) backing length

def putVarint (x : BitVec 64) (length : Nat) (backing : List (BitVec 8)) :
    Option Nat × List (BitVec 8) :=
  putPayload (uvarintBytes (signedVarintWord x)) backing length

theorem putUvarint_result (x : BitVec 64) (length : Nat) (backing : List (BitVec 8)) :
    (putUvarint x length backing).1 =
      if (uvarintBytes x.toNat).length ≤ length then some (uvarintBytes x.toNat).length else none :=
  rfl

theorem putVarint_result (x : BitVec 64) (length : Nat) (backing : List (BitVec 8)) :
    (putVarint x length backing).1 =
      if (uvarintBytes (signedVarintWord x)).length ≤ length
      then some (uvarintBytes (signedVarintWord x)).length else none :=
  rfl

theorem putUvarint_backing (x : BitVec 64) (length : Nat) (backing : List (BitVec 8)) :
    (putUvarint x length backing).2 =
      (uvarintBytes x.toNat).take length ++
      backing.drop (min length (uvarintBytes x.toNat).length) :=
  rfl

theorem putVarint_backing (x : BitVec 64) (length : Nat) (backing : List (BitVec 8)) :
    (putVarint x length backing).2 =
      (uvarintBytes (signedVarintWord x)).take length ++
      backing.drop (min length (uvarintBytes (signedVarintWord x)).length) :=
  rfl

end GoLean.Encoding.Binary
