import GoLean.Image.Color.NRGBA
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [r, g, b, a] := fields | throw (IO.userError "expected four NRGBA channels")
  let (red, green, blue, alpha) := GoLean.Image.Color.nrgbaRGBA
    (← wordInput 8 r) (← wordInput 8 g) (← wordInput 8 b) (← wordInput 8 a)
  return s!"{red} {green} {blue} {alpha}"
