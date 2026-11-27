import GoLean.Compress.Flate.EncoderInit
import GoDiff.Math.Bits.Tuple

open GoLean.Compress.Flate

def main : IO Unit := runTuple fun row => do
  let [width, size] := row | throw (IO.userError "expected native width and size")
  let geometry := match newHuffmanEncoder width.toNat size with
    | none => "bounds"
    | some result => s!"{result.length} {result.capacity} {result.zeroWords}"
  return s!"{geometry}|{maxNode.1.toNat} {maxNode.2.toNat}"
