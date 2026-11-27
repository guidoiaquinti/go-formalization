import GoLean.Compress.Flate.HashWords

namespace GoLean.Compress.Flate

private def loadFour : List (BitVec 8) → Option (BitVec 32)
  | a :: b :: c :: d :: _ => some (d ++ c ++ b ++ a)
  | _ => none

/-- The two-index Go slice checks the backing capacity, even for a short view. -/
def hash4 (backing : List (BitVec 8)) : Option (BitVec 32) :=
  (loadFour backing).map (hash4u · 17)

private def rolling (word : BitVec 32) : List (BitVec 8) → List (BitVec 32)
  | [] => [hash4u word 17]
  | next :: rest =>
    hash4u word 17 :: rolling ((word >>> 8) ||| (next.zeroExtend 32 <<< 24)) rest

/-- The failure flag records a write outside the visible destination. -/
def bulkHash4 (visible : List (BitVec 8)) (length : Nat) (backing : List (BitVec 32)) :
    Bool × List (BitVec 32) :=
  match loadFour visible with
  | none => (false, backing)
  | some first =>
    let results := rolling first (visible.drop 4)
    let written := results.take length
    (length < results.length, written ++ backing.drop written.length)

theorem hash4_ignoreSuffix (a b c d : BitVec 8) (rest : List (BitVec 8)) :
    hash4 (a :: b :: c :: d :: rest) = some (hash4u (d ++ c ++ b ++ a) 17) := rfl

theorem bulkHash4_three (a b c : BitVec 8) (n : Nat) (dst : List (BitVec 32)) :
    bulkHash4 [a, b, c] n dst = (false, dst) := rfl

theorem bulkHash4_first (a b c d : BitVec 8) (old : BitVec 32) (tail : List (BitVec 32)) :
    bulkHash4 [a, b, c, d] 1 (old :: tail) =
      (false, hash4u (d ++ c ++ b ++ a) 17 :: tail) := rfl

end GoLean.Compress.Flate
