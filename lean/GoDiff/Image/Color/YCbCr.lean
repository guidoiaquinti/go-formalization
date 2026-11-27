import GoLean.Image.Color.YCbCr
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [y, cb, cr] := fields | throw (IO.userError "expected three YCbCr fields")
  let (red, green, blue, alpha) :=
    GoLean.Image.Color.yCbCrRGBA (← wordInput 8 y) (← wordInput 8 cb) (← wordInput 8 cr)
  return s!"{red} {green} {blue} {alpha}"
