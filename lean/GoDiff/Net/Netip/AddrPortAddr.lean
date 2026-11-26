import GoLean.Net.Netip.AddrPortAddr
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Netip

def main : IO Unit := runTuple fun fields => do
  let [family, hi, lo, zone, port] := fields | throw (IO.userError "expected family hi lo zone port")
  let detail ← match family with
    | 0 => pure AddrDetailClass.invalid
    | 4 => pure AddrDetailClass.ipv4
    | 6 => pure AddrDetailClass.ipv6
    | _ => throw (IO.userError "expected address class 0, 4 or 6")
  let p := addrPortFrom ⟨detail, UInt64.ofBitVec (← wordInput 64 hi),
    UInt64.ofBitVec (← wordInput 64 lo), if zone = 0 then "" else toString zone⟩
    (UInt16.ofBitVec (← wordInput 16 port))
  let ip := addrPortAddr p
  return s!"{addrBitLen ip.detail} {ip.hi} {ip.lo} [{ip.zone}]"
