import GoLean.Image.Color.NYCbCrAModel
import GoDiff.Math.Bits.Tuple

open GoLean.Image.Color

def main : IO Unit := runTuple fun fields => do
  let [tag, r, g, b, a] := fields |
    throw (IO.userError "expected a type and four channel words")
  let input ← match tag with
    | 0 => pure (NYCbCrAInput.nYCbCrA (r, g, b, a.toNat))
    | 1 => pure (NYCbCrAInput.yCbCr (r, g, b))
    | 2 => pure (NYCbCrAInput.other (some (← wordInput 32 r, ← wordInput 32 g,
        ← wordInput 32 b, ← wordInput 32 a)))
    | 3 => pure (NYCbCrAInput.other none)
    | 4 => pure NYCbCrAInput.nilColor
    | _ => throw (IO.userError "unknown input type")
  let (out, calls) := nYCbCrAModel input
  match out with
  | some (y, cb, cr, alpha) => return s!"{y} {cb} {cr} {alpha} {calls}"
  | none => return s!"panic {calls}"
