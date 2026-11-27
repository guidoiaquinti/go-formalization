import GoLean.Time.ClockArithmetic
import GoDiff.Math.Bits.Tuple

open GoLean.Time

def main : IO Unit := runTuple fun fields => do
  let [tw, te, tl, uw, ue, ul, argument] := fields |
    throw (IO.userError "expected two time states and a duration")
  let t : LocatedTime := ⟨⟨← wordInput 64 tw, BitVec.ofInt 64 te⟩, tl.toNat⟩
  let u : LocatedTime := ⟨⟨← wordInput 64 uw, BitVec.ofInt 64 ue⟩, ul.toNat⟩
  let added := clockAdd t (BitVec.ofInt 64 argument)
  return s!"{added.words.wall.toNat} {added.words.ext.toInt} {added.loc} {(clockSub t u).toInt} {(subMono t.words.ext u.words.ext).toInt}"
