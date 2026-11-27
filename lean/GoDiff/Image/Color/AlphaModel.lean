import GoLean.Image.Color.AlphaModel
import GoDiff.Math.Bits.Tuple

open GoLean.Image.Color

def main : IO Unit := runTuple fun fields => do
  let [tag, value] := fields | throw (IO.userError "expected type and alpha word")
  let input ← match tag with
    | 0 => pure (AlphaInput.alpha value.toNat)
    | 1 => pure (AlphaInput.other (some value.toNat))
    | 2 => pure (AlphaInput.other none)
    | 3 => pure AlphaInput.nilColor
    | _ => throw (IO.userError "unknown input type")
  let (out, calls) := alphaModel input
  return s!"{(out.map toString).getD "panic"} {calls}"
