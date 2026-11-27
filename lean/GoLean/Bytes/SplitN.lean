import GoLean.Strings.SplitN

namespace GoLean.Bytes

/-- Go 1.27.0 bytes.genSplit and bytes.explode cap fragments created by a
    separator or decoded rune. The final raw suffix retains spare capacity. -/
def withSplitCaps (input separator : List Nat) (spare : Nat) (limit : Int)
    (result : Option (List (List Nat))) :
    Option (List (List Nat × Nat)) :=
  let explodeSlots := if limit < 0 then input.length else min limit.toNat input.length
  result.map fun parts =>
    parts.mapIdx fun i part =>
      let rawSuffix := i + 1 == parts.length &&
        (!separator.isEmpty || parts.length == explodeSlots)
      (part, part.length + if rawSuffix then spare else 0)

def splitN (input separator : List Nat) (spare : Nat) (limit : Int) :
    Option (List (List Nat × Nat)) :=
  withSplitCaps input separator spare limit (GoLean.Strings.splitN input separator limit)

theorem splitN_zero (input separator : List Nat) (spare : Nat) :
    splitN input separator spare 0 = none := by
  simp [splitN, withSplitCaps, GoLean.Strings.splitN_zero]

theorem splitN_one (lead : Nat) (rest separator : List Nat) (spare : Nat) :
    splitN (lead :: rest) separator spare 1 =
      some [(lead :: rest, (lead :: rest).length + spare)] := by
  simp [splitN, withSplitCaps, GoLean.Strings.splitN_one, Nat.min_eq_left (Nat.succ_le_succ (Nat.zero_le _))]

end GoLean.Bytes
