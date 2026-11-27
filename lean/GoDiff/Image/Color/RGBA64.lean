import GoLean.Image.Color.RGBA64
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [r, g, b, a] := fields | throw (IO.userError "expected four RGBA channels")
  let (red, green, blue, alpha) := GoLean.Image.Color.rgba64RGBA
    (← wordInput 16 r) (← wordInput 16 g) (← wordInput 16 b) (← wordInput 16 a)
  return s!"{red.toNat} {green.toNat} {blue.toNat} {alpha.toNat}"
