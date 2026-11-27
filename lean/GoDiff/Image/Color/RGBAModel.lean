import GoLean.Image.Color.RGBAModel
import GoDiff.Math.Bits.Tuple

open GoLean.Image.Color

def main : IO Unit := runTuple fun fields => do
  let [width, tag, r, g, b, a] := fields | throw (IO.userError "expected width, tag and four channels")
  let value : Channels := (r.toNat, g.toNat, b.toNat, a.toNat)
  let input ← match tag with
    | 0 => pure (RGBAInput.sameType value)
    | 1 => pure (RGBAInput.other (some value))
    | 2 => pure (RGBAInput.other none)
    | 3 => pure RGBAInput.nilColor
    | _ => throw (IO.userError "unknown input type")
  let (out, calls) := (if width = 8 then rgbaModel else rgba64Model) input
  return match out with
    | some (r, g, b, a) => s!"{r} {g} {b} {a} {calls}"
    | none => s!"panic {calls}"
