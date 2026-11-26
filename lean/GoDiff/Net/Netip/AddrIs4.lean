import GoLean.Net.Netip.AddrIs4
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Netip

def main : IO Unit := runTuple fun fields => do
  let detail ← match fields with
    | [0] => pure AddrDetailClass.invalid
    | [4] => pure AddrDetailClass.ipv4
    | [6] => pure AddrDetailClass.ipv6
    | _ => throw (IO.userError "expected address class 0, 4 or 6")
  return toString (addrIs4 detail)
