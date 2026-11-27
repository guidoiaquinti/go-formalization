import GoLean.Image.Color.Alpha16Model
import GoDiff.Math.Bits.Tuple

open GoLean.Image.Color

def main : IO Unit := runTuple fun fields => do
  let [tag, value] := fields | throw (IO.userError "expected type and alpha word")
  let input ← match tag with
    | 0 => pure (Alpha16Input.alpha16 value.toNat)
    | 1 => pure (Alpha16Input.other (some value.toNat))
    | 2 => pure (Alpha16Input.other none)
    | 3 => pure Alpha16Input.nilColor
    | _ => throw (IO.userError "unknown input type")
  let (out, calls) := alpha16Model input
  return s!"{(out.map toString).getD "panic"} {calls}"
