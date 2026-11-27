import GoLean.Image.Color.YCbCrModel
import GoDiff.Math.Bits.Tuple

open GoLean.Image.Color

def main : IO Unit := runTuple fun fields => do
  let [tag, r, g, b, a] := fields |
    throw (IO.userError "expected a type and four channel words")
  let input ← match tag with
    | 0 => pure (YCbCrInput.yCbCr (r, g, b))
    | 1 => pure (YCbCrInput.other (some (← wordInput 32 r, ← wordInput 32 g,
        ← wordInput 32 b, ← wordInput 32 a)))
    | 2 => pure (YCbCrInput.other none)
    | 3 => pure YCbCrInput.nilColor
    | _ => throw (IO.userError "unknown input type")
  let (out, calls) := yCbCrModel input
  match out with
  | some (y, cb, cr) => return s!"{y} {cb} {cr} {calls}"
  | none => return s!"panic {calls}"
