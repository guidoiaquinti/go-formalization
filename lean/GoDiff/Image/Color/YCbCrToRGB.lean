import GoLean.Image.Color.YCbCrToRGB
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [y, cb, cr] := fields | throw (IO.userError "expected Y, Cb, and Cr")
  let (r, g, b) := GoLean.Image.Color.yCbCrToRGB
    (← wordInput 8 y) (← wordInput 8 cb) (← wordInput 8 cr)
  return s!"{r} {g} {b}"
