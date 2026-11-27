import Init.Data.BitVec.Lemmas

namespace GoLean.Compress.Flate

/-- The output shift is computed in the unsigned byte type. -/
def hash4u (word : BitVec 32) (bits : BitVec 8) : BitVec 32 :=
  (word * 2654435761) >>> ((32 : BitVec 8) - bits).toNat

private def hashWide (word : BitVec 64) (bits : BitVec 8) (shift : Nat)
    (prime : BitVec 64) : BitVec 32 :=
  ((((word <<< shift) * prime) >>> ((64 : BitVec 8) - bits).toNat).truncate 32)

def hashLen (word : BitVec 64) (bits count : BitVec 8) : BitVec 32 :=
  match count.toNat with
  | 3 => (((word <<< (8 : Nat)).truncate 32) * 506832829) >>> ((32 : BitVec 8) - bits).toNat
  | 5 => hashWide word bits 24 889523592379
  | 6 => hashWide word bits 16 227718039650203
  | 7 => hashWide word bits 8 58295818150454627
  | 8 => hashWide word bits 0 0xcf1bbcdcb7a56463
  | _ => hash4u (word.truncate 32) bits

theorem hash4u_full (word : BitVec 32) : hash4u word 32 = word * 2654435761 := by
  simp [hash4u]

theorem hashLen_four (word : BitVec 64) (bits : BitVec 8) :
    hashLen word bits 4 = hash4u (word.truncate 32) bits := rfl

theorem hashLen_default (word : BitVec 64) (bits : BitVec 8) :
    hashLen word bits 255 = hash4u (word.truncate 32) bits := rfl

end GoLean.Compress.Flate
