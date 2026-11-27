import GoLean.Net.Netip.AddrText

namespace GoLean.Net.Netip

private def portTextChunks (p : BinaryPort) : List (List UInt8) :=
  if p.ip.detail = .invalid then []
  else
    let base := addrTextChunks p.ip
    let host := if p.ip.detail = .ipv6 then [[91]] ++ base ++ [[93]] else base
    host ++ [[58], textBytes (toString p.port.toNat)]

def addrPortAppendTo (p : BinaryPort) (b : BinaryBuffer) : BinaryAppendResult :=
  appendTextChunks b (portTextChunks p)

def addrPortAppendText := addrPortAppendTo

def addrPortMarshalText (p : BinaryPort) : List UInt8 := (portTextChunks p).flatten

def addrPortString (p : BinaryPort) : List UInt8 :=
  if p.ip.detail = .invalid then textBytes "invalid AddrPort" else addrPortMarshalText p

private def prefixTextChunks (p : BinaryPrefix) : List (List UInt8) :=
  if p.ip.detail = .invalid then []
  else if p.bits < 0 then [textBytes "invalid Prefix"]
  else addrTextChunks (addrWithoutZone p.ip) ++ [[47]] ++
    (decimalDigits p.bits.toNat.toUInt8).map (fun b => [b])

def prefixAppendTo (p : BinaryPrefix) (b : BinaryBuffer) : BinaryAppendResult :=
  appendTextChunks b (prefixTextChunks p)

def prefixAppendText := prefixAppendTo

def prefixMarshalText (p : BinaryPrefix) : List UInt8 := (prefixTextChunks p).flatten

def prefixString (p : BinaryPrefix) : List UInt8 :=
  if p.bits < 0 then textBytes "invalid Prefix" else prefixMarshalText p

theorem addrPortAppendTo_zero (b : BinaryBuffer) (port : UInt16) :
    addrPortAppendTo ⟨⟨.invalid, 0, 0, []⟩, port⟩ b = appendTextChunks b [] := rfl
theorem addrPortAppendText_same (p : BinaryPort) (b : BinaryBuffer) :
    addrPortAppendText p b = addrPortAppendTo p b := rfl
theorem addrPortMarshalText_zero (port : UInt16) :
    addrPortMarshalText ⟨⟨.invalid, 0, 0, []⟩, port⟩ = [] := rfl
theorem addrPortString_zero (port : UInt16) :
    addrPortString ⟨⟨.invalid, 0, 0, []⟩, port⟩ = textBytes "invalid AddrPort" := rfl
theorem prefixAppendTo_zero (b : BinaryBuffer) :
    prefixAppendTo ⟨⟨.invalid, 0, 0, []⟩, -1⟩ b = appendTextChunks b [] := rfl
theorem prefixAppendText_same (p : BinaryPrefix) (b : BinaryBuffer) :
    prefixAppendText p b = prefixAppendTo p b := rfl
theorem prefixMarshalText_invalidIPv4 (hi lo : UInt64) :
    prefixMarshalText ⟨⟨.ipv4, hi, lo, []⟩, -1⟩ = textBytes "invalid Prefix" := rfl
theorem prefixString_invalid (ip : AddressFields) :
    prefixString ⟨ip, -1⟩ = textBytes "invalid Prefix" := rfl

end GoLean.Net.Netip
