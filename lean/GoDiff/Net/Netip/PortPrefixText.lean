import GoLean.Net.Netip.PortPrefixText
import GoDiff.Math.Bits.Tuple
open GoLean.Net.Netip

private def csv (xs : List UInt8) := String.intercalate "," (xs.map (toString ∘ UInt8.toNat))
private def wire (a : BinaryAppendResult) : String :=
  s!"{a.reused}:{a.nilSlice}:{csv a.output}:{csv a.original}"

def main : IO Unit := runTuple fun row => do
  let f :: hi :: lo :: port :: bits :: nz :: length :: missing :: raw := row |
    throw (IO.userError "expected address, port and prefix fields, then append geometry")
  let cells ← raw.mapM fun x => UInt8.ofBitVec <$> wordInput 8 x
  let ip : AddressFields := ⟨if f = 4 then .ipv4 else if f = 6 then .ipv6 else .invalid,
    UInt64.ofNat hi.toNat, UInt64.ofNat lo.toNat, cells.take nz.toNat⟩
  let b : BinaryBuffer := ⟨length.toNat, cells.drop nz.toNat, missing != 0⟩
  let p : BinaryPort := ⟨ip, UInt16.ofNat port.toNat⟩
  let q : BinaryPrefix := ⟨addrWithoutZone ip, bits⟩
  return String.intercalate "|" [csv (addrPortString p), csv (addrPortMarshalText p),
    wire (addrPortAppendTo p b), wire (addrPortAppendText p b), csv (prefixString q),
    csv (prefixMarshalText q), wire (prefixAppendTo q b), wire (prefixAppendText q b)]
