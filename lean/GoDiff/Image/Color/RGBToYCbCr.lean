import GoLean.Image.Color.RGBToYCbCr
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [r, g, b] := fields | throw (IO.userError "expected red, green and blue bytes")
  let (y, cb, cr) := GoLean.Image.Color.rgbToYCbCr
    (← wordInput 8 r) (← wordInput 8 g) (← wordInput 8 b)
  return s!"{y} {cb} {cr}"
