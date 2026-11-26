import GoLean.Net.Netip.PrefixIsSingleIP
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Netip

def main : IO Unit := runTuple fun fields => do
  let [width, bits] := fields | throw (IO.userError "expected width and bits")
  let detail ← match width with
    | 0 => pure AddrDetailClass.invalid
    | 32 => pure AddrDetailClass.ipv4
    | 128 => pure AddrDetailClass.ipv6
    | _ => throw (IO.userError "invalid address width")
  return toString (prefixIsSingleIP (prefixFrom ⟨detail, 0, 0, ""⟩ bits))
