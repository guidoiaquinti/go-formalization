import GoLean.Image.Color.RGBToCMYK
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [r, g, b] := fields | throw (IO.userError "expected three RGB channels")
  let (cyan, magenta, yellow, black) := GoLean.Image.Color.rgbToCMYK
    (← wordInput 8 r) (← wordInput 8 g) (← wordInput 8 b)
  return s!"{cyan.toNat} {magenta.toNat} {yellow.toNat} {black.toNat}"
