import GoLean.Net.Netip.AddrFrom4
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  unless fields.all (fun x => 0 ≤ x && x < 256) do
    throw (IO.userError "expected bytes")
  let bytes := fields.toArray.map (fun x => x.toNat.toUInt8)
  if h : bytes.size = 4 then
    let addr := GoLean.Net.Netip.addrFrom4 fun i => bytes[i.val]'(by omega)
    return s!"{addr.hi} {addr.lo} {addr.isV6} [{addr.zone}]"
  else
    throw (IO.userError "expected 4 bytes")
