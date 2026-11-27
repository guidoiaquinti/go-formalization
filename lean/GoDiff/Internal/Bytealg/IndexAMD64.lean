import GoLean.Internal.Bytealg.IndexAMD64
import GoDiff.Math.Bits.Tuple

open GoLean.Internal.Bytealg

def main : IO Unit := runTuple fun row => do
  let [count] := row | throw (IO.userError "expected signed processed count")
  return s!"{amd64Cutover (BitVec.ofInt 64 count)},{amd64IndexInit false},{amd64IndexInit true}"
