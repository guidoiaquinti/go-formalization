import GoLean.Image.Color.CMYKModel
import GoDiff.Math.Bits.Tuple

open GoLean.Image.Color

def main : IO Unit := runTuple fun fields => do
  let [tag, r, g, b, a] := fields |
    throw (IO.userError "expected a type and four channel words")
  let input ← match tag with
    | 0 => pure (CMYKInput.cmyk (← wordInput 8 r, ← wordInput 8 g,
        ← wordInput 8 b, ← wordInput 8 a))
    | 1 => pure (CMYKInput.other (some (← wordInput 32 r, ← wordInput 32 g,
        ← wordInput 32 b, ← wordInput 32 a)))
    | 2 => pure (CMYKInput.other none)
    | 3 => pure CMYKInput.nilColor
    | _ => throw (IO.userError "unknown input type")
  let (out, calls) := cmykModel input
  match out with
  | some (c, m, y, k) => return s!"{c.toNat} {m.toNat} {y.toNat} {k.toNat} {calls}"
  | none => return s!"panic {calls}"
