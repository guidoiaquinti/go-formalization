import GoLean.Net.Netip.TextDigits
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Netip

private def csv (xs : List UInt8) := String.intercalate "," (xs.map (toString ∘ UInt8.toNat))
private def wire (a : BinaryAppendResult) : String :=
  s!"{a.reused}:{csv a.output}:{csv a.original}"

def main : IO Unit := runTuple fun row => do
  let value :: len :: raw := row | throw (IO.userError "expected value, visible length and backing")
  let word ← wordInput 16 value
  let input ← raw.mapM fun x => UInt8.ofBitVec <$> wordInput 8 x
  let b : BinaryBuffer := ⟨len.toNat, input, input.isEmpty⟩
  return String.intercalate "|" [wire (appendDecimal b value.toNat.toUInt8),
    wire (appendHex b (UInt16.ofBitVec word)), wire (appendHexPad b (UInt16.ofBitVec word))]
