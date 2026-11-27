import GoLean.Time.Location
import GoDiff.Math.Bits.Tuple

open GoLean.Time

private def wire (s : LocatedTime) := s!"{s.words.wall.toNat},{s.words.ext.toNat},{s.loc}"

def main : IO Unit := runTuple fun row => do
  let [wall, ext, old, next] := row | throw (IO.userError "expected clock words and location tags")
  let wall ← wordInput 64 wall
  let ext ← wordInput 64 ext
  let value : LocatedTime := ⟨⟨wall, ext⟩, old.toNat⟩
  return String.intercalate "|" [wire (clockStripMono value), wire (clockSetLoc value next.toNat),
    wire (clockUTC value), wire (clockLocal value), (clockIn value next.toNat).map wire |>.getD "nil",
    toString (clockLocation value)]
