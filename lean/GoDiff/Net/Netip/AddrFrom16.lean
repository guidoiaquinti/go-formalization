import GoLean.Net.Netip.AddrFrom16
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  unless fields.all (fun x => 0 ≤ x && x < 256) do
    throw (IO.userError "expected bytes")
  let bytes := fields.toArray.map (fun x => x.toNat.toUInt8)
  if h : bytes.size = 16 then
    let addr := GoLean.Net.Netip.addrFrom16 fun i => bytes[i.val]'(by omega)
    return s!"{addr.hi} {addr.lo} {addr.isV6} [{addr.zone}]"
  else
    throw (IO.userError "expected 16 bytes")
