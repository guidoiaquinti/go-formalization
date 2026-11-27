import GoLean.Strings.MakeASCIISet
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let _ :: raw := row | throw (IO.userError "expected input marker and bytes")
  let bytes ← raw.mapM (wordInput 8)
  let (table, valid) := GoLean.Strings.makeASCIISet bytes
  let entries := (List.range 256).map fun i => if table (BitVec.ofNat 8 i) then '1' else '0'
  return s!"{valid}|{String.ofList entries}"
