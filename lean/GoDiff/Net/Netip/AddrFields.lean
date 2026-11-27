import GoLean.Net.Netip.AddrFields
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Netip

private def csv (xs : List UInt8) := String.intercalate "," (xs.map (toString ∘ UInt8.toNat))

private def wire (ip : AddressFields) : String :=
  let family := match ip.detail with | .invalid => 0 | .ipv4 => 4 | .ipv6 => 6
  s!"{family}:{ip.hi.toNat}:{ip.lo.toNat}:{csv ip.zone}"

def main : IO Unit := runTuple fun row => do
  let family :: hi :: lo :: index :: zoneLength :: raw := row |
    throw (IO.userError "expected family, limbs, index, zone length and raw bytes")
  let detail := if family = 4 then AddrDetailClass.ipv4 else
    if family = 6 then AddrDetailClass.ipv6 else AddrDetailClass.invalid
  let bytes ← raw.mapM fun v => UInt8.ofBitVec <$> wordInput 8 v
  let ip : AddressFields := ⟨detail, UInt64.ofNat hi.toNat, UInt64.ofNat lo.toNat, bytes.take zoneLength.toNat⟩
  let i := UInt8.ofNat index.toNat
  let fromSlice := (addrFromSlice (bytes.drop zoneLength.toNat)).map fun addr =>
    wire ⟨if addr.isV6 then .ipv6 else .ipv4, addr.hi, addr.lo, []⟩
  return s!"{addrV4Byte ip i}|{addrV6Byte ip i}|{addrV6Word ip i}|{addrIsZero ip}|{addrHasZone ip}|{wire (addrWithoutZone ip)}|{fromSlice.getD "none"}"
