import GoLean.Math.Ldexp
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [width, word, exponent] := fields | throw (IO.userError "expected width, binary64 word and exponent")
  if width != 32 && width != 64 then throw (IO.userError "expected a native Go int width")
  let bits ← wordInput 64 word
  return toString (GoLean.Math.ldexp width.toNat bits.toNat exponent)
