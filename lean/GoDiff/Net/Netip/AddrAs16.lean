import GoLean.Net.Netip.AddrAs16
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [hi, lo] := fields | throw (IO.userError "expected hi lo")
  let hi := UInt64.ofBitVec (← wordInput 64 hi)
  let lo := UInt64.ofBitVec (← wordInput 64 lo)
  let bytes := List.ofFn (GoLean.Net.Netip.addrAs16 hi lo)
  return String.intercalate " " (bytes.map toString)
