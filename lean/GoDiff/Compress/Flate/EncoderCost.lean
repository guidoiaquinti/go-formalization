import GoLean.Compress.Flate.EncoderCost
import GoDiff.Math.Bits.Tuple

open GoLean.Compress.Flate

private def wire : CostResult → String
  | .nilReceiver => "nil"
  | .bounds => "bounds"
  | .value n => toString n

def main : IO Unit := runTuple fun row => do
  let width :: present :: size :: frequencies :: values := row |
    throw (IO.userError "expected width, receiver, table length and frequency length")
  let codes ← (values.take size.toNat).mapM (wordInput 32)
  let freqs ← ((values.drop size.toNat).take frequencies.toNat).mapM (wordInput 16)
  let data ← (values.drop (size + frequencies).toNat).mapM (wordInput 8)
  let codes := if present == 0 then none else some codes
  return String.intercalate " " [wire (bitLength width.toNat codes freqs),
    wire (canEncodeLen width.toNat codes freqs), wire (bitLengthRaw width.toNat codes data)]
