import GoLean.Image.Color.RGBA
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [r, g, b, a] := fields | throw (IO.userError "expected four RGBA channels")
  let (red, green, blue, alpha) := GoLean.Image.Color.rgbaRGBA
    (← wordInput 8 r) (← wordInput 8 g) (← wordInput 8 b) (← wordInput 8 a)
  return s!"{red.toNat} {green.toNat} {blue.toNat} {alpha.toNat}"
