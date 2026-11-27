import GoLean.Net.Netip.AddrText
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Netip

private def csv (xs : List UInt8) := String.intercalate "," (xs.map (toString ∘ UInt8.toNat))
private def wire (a : BinaryAppendResult) : String :=
  s!"{a.reused}:{a.nilSlice}:{csv a.output}:{csv a.original}"

def main : IO Unit := runTuple fun row => do
  let family :: hi :: lo :: zlen :: len :: missing :: raw := row |
    throw (IO.userError "expected address, zone and append geometry")
  let bytes ← raw.mapM fun x => UInt8.ofBitVec <$> wordInput 8 x
  let ip : AddressFields := ⟨if family = 4 then .ipv4 else if family = 6 then .ipv6 else .invalid,
    UInt64.ofNat hi.toNat, UInt64.ofNat lo.toNat, bytes.take zlen.toNat⟩
  let b : BinaryBuffer := ⟨len.toNat, bytes.drop zlen.toNat, missing != 0⟩
  let privateOut := if family = 4 then wire (addrAppendTo4 ip b)
    else if family = 6 then
      if hi = 0 && lo.toNat / 2^32 = 65535 then wire (addrAppendTo4In6 ip b)
      else wire (addrAppendTo6 ip b)
    else "skip"
  return String.intercalate "|" [csv (addrString ip), csv (addrStringExpanded ip),
    csv (addrMarshalText ip), wire (addrAppendTo ip b), wire (addrAppendText ip b), privateOut]
