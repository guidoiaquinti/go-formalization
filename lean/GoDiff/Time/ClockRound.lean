import GoLean.Time.ClockRound
import GoDiff.Math.Bits.Tuple

open GoLean.Time

private def wire (s : LocatedTime) := s!"{s.words.wall.toNat},{s.words.ext.toInt},{s.loc}"

def main : IO Unit := runTuple fun row => do
  let [wall, ext, loc, duration, half] := row | throw (IO.userError "expected raw time and duration")
  let wall ← wordInput 64 wall
  let s := ⟨⟨wall, BitVec.ofInt 64 ext⟩, loc.toNat⟩
  let d := BitVec.ofInt 64 duration
  let privateResult := if duration ≤ 0 then "skip" else
    let (q,r) := clockDiv s.words duration
    s!"{q},{r}"
  return String.intercalate "|" [wire (clockTruncate s d), wire (clockRound s d),
    privateResult, toString (lessThanHalf (BitVec.ofInt 64 half) d)]
