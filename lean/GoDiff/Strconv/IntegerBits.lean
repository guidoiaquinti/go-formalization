import GoLean.Strconv.IntegerBits
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let [width, value] := row | throw (IO.userError "expected width and word")
  if width == 8 then
    return toString (GoLean.Strconv.lower (← wordInput 8 value)).toNat
  if width == 32 then
    return toString (GoLean.Strconv.isPowerOfTwo (← wordInput 32 value))
  if width == 64 then
    return toString (GoLean.Strconv.isPowerOfTwo (← wordInput 64 value))
  throw (IO.userError "unsupported word width")
