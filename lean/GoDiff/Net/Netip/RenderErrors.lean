import GoLean.Net.Netip.RenderErrors
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Netip

private def wire (cells : ParseBytes) := String.intercalate "," (cells.map toString)

def main : IO Unit := runTuple fun row => do
  let inputLength :: msgLength :: raw := row | throw (IO.userError "expected two lengths")
  let raw := raw.map Int.toNat
  let input := raw.take inputLength.toNat
  let message := (raw.drop inputLength.toNat).take msgLength.toNat
  let suffix := raw.drop (inputLength + msgLength).toNat
  return s!"{wire (renderAddrError input message suffix)}|{wire (renderPrefixError input message)}"
