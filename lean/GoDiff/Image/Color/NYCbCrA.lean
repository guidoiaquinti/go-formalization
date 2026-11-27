import GoLean.Image.Color.NYCbCrA
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [y, cb, cr, a] := fields | throw (IO.userError "expected Y, Cb, Cr and alpha")
  let (red, green, blue, alpha) := GoLean.Image.Color.nYCbCrARGBA
    (← wordInput 8 y) (← wordInput 8 cb) (← wordInput 8 cr) (← wordInput 8 a)
  return s!"{red} {green} {blue} {alpha}"
