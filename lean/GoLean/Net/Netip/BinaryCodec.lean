import GoLean.Net.Netip.AddrFields

namespace GoLean.Net.Netip

structure BinaryBuffer where
  length : Nat
  cells : List UInt8
  nilSlice : Bool

structure BinaryAppendResult where
  output : List UInt8
  original : List UInt8
  reused : Bool
  nilSlice : Bool

private def bytes (word width : Nat) : List UInt8 :=
  (List.range width).map fun i => UInt8.ofNat (word / 256^(width - 1 - i))

private def addrChunks (ip : AddressFields) : List (List UInt8) :=
  match ip.detail with
  | .invalid => []
  | .ipv4 => [bytes ip.lo.toNat 4]
  | .ipv6 => [bytes ip.hi.toNat 8, bytes ip.lo.toNat 8, ip.zone]

private def appendChunks (buf : BinaryBuffer) (parts : List (List UInt8)) : BinaryAppendResult :=
  let update := fun (old, cursor, reused) (part : List UInt8) =>
    if reused && cursor + part.length ≤ buf.cells.length then
      (old.take cursor ++ part ++ old.drop (cursor + part.length), cursor + part.length, true)
    else (old, cursor, false)
  let (original, _, reused) := parts.foldl update (buf.cells, buf.length, true)
  let output := buf.cells.take buf.length ++ parts.flatten
  ⟨output, original, reused, buf.nilSlice && output.isEmpty⟩

def addrAppendBinary (ip : AddressFields) (buf : BinaryBuffer) : BinaryAppendResult :=
  appendChunks buf (addrChunks ip)

def addrMarshalBinarySize (ip : AddressFields) : Nat :=
  match ip.detail with
  | .invalid => 0
  | .ipv4 => 4
  | .ipv6 => 16 + ip.zone.length

def addrMarshalBinary (ip : AddressFields) : List UInt8 := (addrChunks ip).flatten

private def zero : AddressFields := ⟨.invalid, 0, 0, []⟩

private def unmarshalAddr (raw : List UInt8) : Option AddressFields :=
  if raw.isEmpty then some zero
  else
    let addressBytes := raw.take 16
    (addrFromSlice addressBytes).map fun ip =>
      ⟨if ip.isV6 then .ipv6 else .ipv4, ip.hi, ip.lo, raw.drop 16⟩

def addrUnmarshalBinary (old : AddressFields) (raw : List UInt8) : Bool × AddressFields :=
  match unmarshalAddr raw with
  | none => (false, old)
  | some ip => (true, ip)

structure BinaryPort where
  ip : AddressFields
  port : UInt16

private def portChunks (p : BinaryPort) : List (List UInt8) :=
  addrChunks p.ip ++ [(bytes p.port.toNat 2).reverse]

def portAppendBinary (p : BinaryPort) (buf : BinaryBuffer) : BinaryAppendResult :=
  appendChunks buf (portChunks p)

def portMarshalBinary (p : BinaryPort) : List UInt8 := (portChunks p).flatten

def portUnmarshalBinary (old : BinaryPort) (raw : List UInt8) : Bool × BinaryPort :=
  if raw.length < 2 then (false, old) else
    let count := raw.length - 2
    match unmarshalAddr (raw.take count) with
    | none => (false, old)
    | some ip => (true, ⟨ip, UInt16.ofNat (raw[count]!.toNat + 256 * raw[count + 1]!.toNat)⟩)

structure BinaryPrefix where
  ip : AddressFields
  bits : Int

private def prefixChunks (p : BinaryPrefix) : List (List UInt8) :=
  addrChunks (addrWithoutZone p.ip) ++ [[(BitVec.ofInt 8 p.bits).toNat.toUInt8]]

def prefixAppendBinary (p : BinaryPrefix) (buf : BinaryBuffer) : BinaryAppendResult :=
  appendChunks buf (prefixChunks p)

def prefixMarshalBinary (p : BinaryPrefix) : List UInt8 := (prefixChunks p).flatten

def prefixUnmarshalBinary (old : BinaryPrefix) (raw : List UInt8) : Bool × BinaryPrefix :=
  if raw.isEmpty then (false, old) else
    let count := raw.length - 1
    match unmarshalAddr (raw.take count) with
    | none => (false, old)
    | some ip =>
      let bit := raw[count]!.toNat
      let valid := ip.detail != .invalid && bit ≤ addrBitLen ip.detail
      (true, ⟨addrWithoutZone ip, if valid then Int.ofNat bit else -1⟩)

theorem addrAppendBinary_zero (b : BinaryBuffer) :
    addrAppendBinary zero b = ⟨b.cells.take b.length, b.cells, true,
      b.nilSlice && (b.cells.take b.length).isEmpty⟩ := by simp [addrAppendBinary, addrChunks, zero, appendChunks]
theorem addrMarshalBinarySize_zero : addrMarshalBinarySize zero = 0 := rfl
theorem addrMarshalBinary_zero : addrMarshalBinary zero = [] := rfl
theorem addrUnmarshalBinary_empty (old : AddressFields) : addrUnmarshalBinary old [] = (true, zero) := rfl
theorem portAppendBinary_zero (b : BinaryBuffer) :
    portAppendBinary ⟨zero, 0⟩ b = appendChunks b [[0, 0]] := rfl
theorem portMarshalBinary_zero : portMarshalBinary ⟨zero, 0⟩ = [0, 0] := rfl
theorem portUnmarshalBinary_empty (old : BinaryPort) : portUnmarshalBinary old [] = (false, old) := rfl
theorem prefixAppendBinary_zero (b : BinaryBuffer) :
    prefixAppendBinary ⟨zero, -1⟩ b = appendChunks b [[255]] := rfl
theorem prefixMarshalBinary_zero : prefixMarshalBinary ⟨zero, -1⟩ = [255] := rfl
theorem prefixUnmarshalBinary_empty (old : BinaryPrefix) : prefixUnmarshalBinary old [] = (false, old) := rfl

end GoLean.Net.Netip
