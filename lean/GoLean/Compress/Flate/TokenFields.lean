import Init.Data.BitVec.Lemmas
import Init.Omega

namespace GoLean.Compress.Flate

def tokenType (word : BitVec 32) : Nat := word.toNat / 1073741824 * 1073741824

def tokenLiteral (word : BitVec 32) : Nat := word.toNat % 256

def tokenOffset (word : BitVec 32) : Nat := word.toNat % 4194304

def tokenLength (word : BitVec 32) : Nat := word.toNat / 4194304 % 256

theorem tokenFields_partition (word : BitVec 32) :
    tokenType word + tokenLength word * 4194304 + tokenOffset word = word.toNat := by
  unfold tokenType tokenLength tokenOffset
  omega

theorem tokenLiteral_offset (word : BitVec 32) :
    tokenLiteral word = tokenOffset word % 256 := by
  unfold tokenLiteral tokenOffset
  omega

end GoLean.Compress.Flate
