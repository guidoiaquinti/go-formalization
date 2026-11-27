import GoLean.Image.Color.Gray16Model
import GoDiff.Math.Bits.Tuple

open GoLean.Image.Color

def main : IO Unit := runTuple fun fields => do
  let [tag, r, g, b, _] := fields | throw (IO.userError "expected type and RGBA words")
  let input ← match tag with
    | 0 => pure (GrayInput.gray r.toNat)
    | 1 => pure (GrayInput.other (some (r.toNat, g.toNat, b.toNat)))
    | 2 => pure (GrayInput.other none)
    | 3 => pure GrayInput.nilColor
    | _ => throw (IO.userError "unknown input type")
  let (out, calls) := gray16Model input
  return s!"{(out.map toString).getD "panic"} {calls}"
