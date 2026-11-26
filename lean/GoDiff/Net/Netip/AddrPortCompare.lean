import GoLean.Net.Netip.AddrPortCompare
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Netip

private def address (family hi lo zone : Int) : IO ZonedAddr := do
  let detail ← match family with
    | 0 => pure AddrDetailClass.invalid
    | 4 => pure AddrDetailClass.ipv4
    | 6 => pure AddrDetailClass.ipv6
    | _ => throw (IO.userError "expected address class 0, 4 or 6")
  return ⟨detail, UInt64.ofBitVec (← wordInput 64 hi),
    UInt64.ofBitVec (← wordInput 64 lo), if zone = 0 then "" else toString zone⟩

def main : IO Unit := runTuple fun fields => do
  let [af, ah, al, az, ap, bf, bh, bl, bz, bp] := fields | throw (IO.userError "expected two addresses")
  let a := addrPortFrom (← address af ah al az) (UInt16.ofBitVec (← wordInput 16 ap))
  let b := addrPortFrom (← address bf bh bl bz) (UInt16.ofBitVec (← wordInput 16 bp))
  return toString (addrPortCompare a b)
