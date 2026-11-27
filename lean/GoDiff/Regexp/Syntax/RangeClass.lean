import GoLean.Regexp.Syntax.RangeClass
import GoDiff.Math.Bits.Tuple

open GoLean.Regexp.Syntax

private def pairs : List Int → List RuneRange
  | lo :: hi :: rest => (lo.toNat, hi.toNat) :: pairs rest
  | _ => []

private def encode (ranges : List RuneRange) : String :=
  String.intercalate "," (ranges.flatMap fun (lo,hi) => [toString lo, toString hi])

private def wire (store : ClassStore) : String :=
  s!"{encode store.visible}:{encode store.backing}:{store.shared}"

def main : IO Unit := runTuple fun fields => do
  let size :: capacity :: low :: high :: rune :: i :: j :: raw := fields |
    throw (IO.userError "expected range geometry and endpoints")
  let entries := pairs raw
  let backing := entries.take capacity.toNat
  let visible := backing.take size.toNat
  let clean := entries.drop capacity.toNat
  let t : ClassStore := ⟨visible, backing, true⟩
  let compare := if visible.isEmpty then "skip" else
    s!"{rangesLess visible i.toNat j.toNat}:{encode (rangesSwap visible i.toNat j.toNat ++ backing.drop size.toNat)}"
  let neg : ClassStore := ⟨clean, clean, true⟩
  return String.intercalate "|" [toString (rangesLen visible), compare, wire (appendRange t (low.toNat, high.toNat)),
    wire (appendClass t clean), wire (appendNegatedClass t clean), wire (negateClass neg),
    toString (inCharClass rune.toNat clean)]
