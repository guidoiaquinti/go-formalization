import GoLean.Net.Netip.AddrAs4
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Netip

def main : IO Unit := runTuple fun fields => do
  let [family, hi, lo] := fields | throw (IO.userError "expected family hi lo")
  let detail ← match family with
    | 0 => pure AddrDetailClass.invalid
    | 4 => pure AddrDetailClass.ipv4
    | 6 => pure AddrDetailClass.ipv6
    | _ => throw (IO.userError "expected address class 0, 4 or 6")
  match addrAs4 detail (UInt64.ofBitVec (← wordInput 64 hi)) (UInt64.ofBitVec (← wordInput 64 lo)) with
  | .error msg => return "panic: " ++ msg
  | .ok bytes => return String.intercalate " " ((List.ofFn bytes).map toString)
