import GoLean.Compress.Flate.FixedHuffman
import GoDiff.Math.Bits.Tuple

open GoLean.Compress.Flate

def main : IO Unit := runTuple fun row => do
  let [kind] := row | throw (IO.userError "expected the fixed table kind")
  let table ← if kind = 0 then pure generateFixedLiteralEncoding
    else if kind = 1 then pure generateFixedOffsetEncoding
    else throw (IO.userError "invalid table kind")
  let tail := List.replicate (table.capacity - table.codes.length) (0 : HCode)
  let words := " ".intercalate ((table.codes ++ tail).map (toString ∘ BitVec.toNat))
  return s!"{table.codes.length} {table.capacity} {table.zeroWorkspace} {words}"
