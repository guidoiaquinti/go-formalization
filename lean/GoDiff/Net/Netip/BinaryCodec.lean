import GoLean.Net.Netip.BinaryCodec
import GoDiff.Math.Bits.Tuple
open GoLean.Net.Netip

private def csv (xs : List UInt8) := String.intercalate "," (xs.map (toString ∘ UInt8.toNat))

private def ipWire (ip : AddressFields) : String :=
  let family := match ip.detail with | .invalid => 0 | .ipv4 => 4 | .ipv6 => 6
  s!"{family}:{ip.hi.toNat}:{ip.lo.toNat}:{csv ip.zone}"

private def bufferWire (a : BinaryAppendResult) : String :=
  s!"{a.reused}:{a.nilSlice}:{csv a.output}:{csv a.original}"

def main : IO Unit := runTuple fun row => do
  let f :: hi :: lo :: port :: bits :: nz :: len :: cap :: noBuf :: raw := row |
    throw (IO.userError "expected address, metadata, backing and wire bytes")
  let detail := if f = 4 then AddrDetailClass.ipv4 else if f = 6 then .ipv6 else .invalid
  let cells ← raw.mapM fun b => UInt8.ofBitVec <$> wordInput 8 b
  let rest := cells.drop nz.toNat
  let ip : AddressFields := ⟨detail, UInt64.ofNat hi.toNat, UInt64.ofNat lo.toNat, cells.take nz.toNat⟩
  let buf : BinaryBuffer := ⟨len.toNat, rest.take cap.toNat, noBuf != 0⟩
  let payload := rest.drop cap.toNat
  let p : BinaryPort := ⟨ip, UInt16.ofNat port.toNat⟩
  let q : BinaryPrefix := ⟨addrWithoutZone ip, bits⟩
  let (okA, a) := addrUnmarshalBinary ip payload
  let (okP, p2) := portUnmarshalBinary p payload
  let (okQ, q2) := prefixUnmarshalBinary q payload
  return String.intercalate "|" [
    toString (addrMarshalBinarySize ip), csv (addrMarshalBinary ip), bufferWire (addrAppendBinary ip buf),
    s!"{okA}:{ipWire a}", csv (portMarshalBinary p), bufferWire (portAppendBinary p buf),
    s!"{okP}:{ipWire p2.ip}:{p2.port}", csv (prefixMarshalBinary q), bufferWire (prefixAppendBinary q buf),
    s!"{okQ}:{ipWire q2.ip}:{q2.bits}"]
