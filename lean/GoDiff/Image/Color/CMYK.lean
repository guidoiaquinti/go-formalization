import GoLean.Image.Color.CMYK
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [c, m, y, k] := fields | throw (IO.userError "expected four CMYK channels")
  let (red, green, blue, alpha) := GoLean.Image.Color.cmykRGBA
    (← wordInput 8 c) (← wordInput 8 m) (← wordInput 8 y) (← wordInput 8 k)
  return s!"{red.toNat} {green.toNat} {blue.toNat} {alpha.toNat}"
