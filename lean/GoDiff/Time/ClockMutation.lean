import GoLean.Time.ClockMutation
import GoDiff.Math.Bits.Tuple

open GoLean.Time

private def wire (s : LocatedTime) := s!"{s.words.wall.toNat},{s.words.ext.toNat},{s.loc}"

def main : IO Unit := runTuple fun fields => do
  let [wall, ext, loc, argument] := fields | throw (IO.userError "expected raw clock and argument")
  let wall ← wordInput 64 wall
  let ext ← wordInput 64 ext
  let argument ← wordInput 64 argument
  let state : LocatedTime := ⟨⟨wall, ext⟩, loc.toNat⟩
  return s!"{wire (clockSetMono state argument)}|{wire (clockAddSec state argument)}"
