import GoLean.Math.Logb
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [word] := fields | throw (IO.userError "expected one binary64 word")
  let raw ← wordInput 64 word
  return match GoLean.Math.logb raw.toNat with
  | .finite exponent => s!"finite {exponent}"
  | .special bits => s!"special {bits}"
