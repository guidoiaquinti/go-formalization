import GoLean.Net.Netip.PrefixAddr
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Netip

def main : IO Unit := runTuple fun fields => do
  let [family, hi, lo, bits] := fields | throw (IO.userError "expected family hi lo bits")
  let detail ← match family with
    | 0 => pure AddrDetailClass.invalid
    | 4 => pure AddrDetailClass.ipv4
    | 6 => pure AddrDetailClass.ipv6
    | _ => throw (IO.userError "expected address class 0, 4 or 6")
  let width := addrBitLen detail
  let bitsPlusOne := if width != 0 && 0 ≤ bits && bits ≤ Int.ofNat width
    then (bits + 1).toNat.toUInt8 else 0
  let ip := prefixAddr ⟨⟨detail, UInt64.ofBitVec (← wordInput 64 hi),
    UInt64.ofBitVec (← wordInput 64 lo), ""⟩, bitsPlusOne⟩
  return s!"{addrBitLen ip.detail} {ip.hi} {ip.lo} [{ip.zone}]"
