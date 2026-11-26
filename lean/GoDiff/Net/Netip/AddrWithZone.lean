import GoLean.Net.Netip.AddrWithZone
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Netip

def main : IO Unit := runTuple fun fields => do
  let [family, hi, lo, zone, requested] := fields | throw (IO.userError "expected family hi lo zone requested")
  let detail ← match family with
    | 0 => pure AddrDetailClass.invalid
    | 4 => pure AddrDetailClass.ipv4
    | 6 => pure AddrDetailClass.ipv6
    | _ => throw (IO.userError "expected address class 0, 4 or 6")
  let ip := addrWithZone ⟨detail, UInt64.ofBitVec (← wordInput 64 hi),
    UInt64.ofBitVec (← wordInput 64 lo), if zone = 0 then "" else toString zone⟩
    (if requested = 0 then "" else toString requested)
  return s!"{addrBitLen ip.detail} {ip.hi} {ip.lo} [{ip.zone}]"
