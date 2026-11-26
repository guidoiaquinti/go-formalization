import GoLean.Net.Netip.AddrIs4In6
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Netip

def main : IO Unit := runTuple fun fields => do
  let [family, hi, lo] := fields | throw (IO.userError "expected family hi lo")
  let detail ← match family with
    | 0 => pure AddrDetailClass.invalid
    | 4 => pure AddrDetailClass.ipv4
    | 6 => pure AddrDetailClass.ipv6
    | _ => throw (IO.userError "expected address class 0, 4 or 6")
  return toString (addrIs4In6 detail (UInt64.ofBitVec (← wordInput 64 hi)) (UInt64.ofBitVec (← wordInput 64 lo)))
