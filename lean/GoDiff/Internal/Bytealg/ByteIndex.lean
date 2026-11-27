import GoLean.Internal.Bytealg.ByteIndex
import GoDiff.Math.Bits.Tuple

open GoLean.Internal.Bytealg

def main : IO Unit := runTuple fun row => do
  let target :: values := row | throw (IO.userError "expected target byte")
  let target := UInt8.ofNat (← wordInput 8 target).toNat
  let bytes ← values.mapM fun value => do pure (UInt8.ofNat (← wordInput 8 value).toNat)
  let indices := [indexByte bytes target, indexByteString bytes target,
    lastIndexByte bytes target, lastIndexByteString bytes target]
  return String.intercalate "," (indices.map fun index => index.map toString |>.getD "-1")
