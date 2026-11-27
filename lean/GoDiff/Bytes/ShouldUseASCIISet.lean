import GoLean.Bytes.ShouldUseASCIISet
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [bufLen] := fields | throw (IO.userError "expected one signed buffer length")
  return toString (GoLean.Bytes.shouldUseASCIISet bufLen)
