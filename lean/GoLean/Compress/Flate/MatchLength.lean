import Init.Data.BitVec.Lemmas

namespace GoLean.Compress.Flate

private def commonBytes : List (BitVec 8) → List (BitVec 8) → Nat
  | a :: as, b :: bs => if a == b then 1 + commonBytes as bs else 0
  | _, _ => 0

private def matchBlocks (rightLength rightCapacity : Nat) :
    Nat → Nat → List (BitVec 8) → List (BitVec 8) → Option Nat
  | 0, consumed, left, right =>
    if rightLength < consumed || rightCapacity < consumed + left.length then none
    else some (consumed + commonBytes left right)
  | fuel + 1, consumed, left, right =>
    if rightCapacity < consumed + 8 then none
    else
      let equal := commonBytes (left.take 8) (right.take 8)
      if equal < 8 then some (consumed + equal)
      else matchBlocks rightLength rightCapacity fuel (consumed + 8) (left.drop 8) (right.drop 8)

/-- Word chunks can read through the second slice's capacity. Entering the
    tail also requires that the consumed prefix fits the second visible view. -/
def matchLen (left : List (BitVec 8)) (rightLength : Nat) (right : List (BitVec 8)) :
    Option Nat :=
  matchBlocks rightLength right.length (left.length / 8) 0 left right

theorem matchLen_empty (n : Nat) (right : List (BitVec 8)) : matchLen [] n right = some 0 := by
  simp [matchLen, matchBlocks, commonBytes]

theorem matchLen_tailCapacity (a b : BitVec 8) (tail : List (BitVec 8)) :
    matchLen [a, b] 0 (a :: b :: tail) = some 2 := by
  simp [matchLen, matchBlocks, commonBytes]

theorem matchLen_blockTailBounds (a : BitVec 8) :
    matchLen (List.replicate 8 a) 7 (List.replicate 8 a) = none := by
  simp [matchLen, matchBlocks, commonBytes, List.replicate_succ]

end GoLean.Compress.Flate
