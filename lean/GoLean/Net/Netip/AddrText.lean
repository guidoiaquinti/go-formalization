import GoLean.Net.Netip.TextDigits

namespace GoLean.Net.Netip

def textBytes (s : String) : List UInt8 := s.toUTF8.data.toList

def appendTextChunks (b : BinaryBuffer) (chunks : List (List UInt8)) : BinaryAppendResult :=
  let (old, _, reused) := chunks.foldl (fun (old, cursor, same) bytes =>
    if same && cursor + bytes.length ≤ b.cells.length then
      (old.take cursor ++ bytes ++ old.drop (cursor + bytes.length), cursor + bytes.length, true)
    else (old, cursor, false)) (b.cells, b.length, true)
  let out := b.cells.take b.length ++ chunks.flatten
  ⟨out, old, reused, b.nilSlice && out.isEmpty⟩

private def zone (a : AddressFields) : List (List UInt8) :=
  if a.zone.isEmpty then [] else [[37], a.zone]

private def ipv4Chunks (a : AddressFields) : List (List UInt8) :=
  let octets := (List.range 4).map fun i =>
    (decimalDigits (addrV4Byte a (UInt8.ofNat i))).map (fun b => [b])
  (octets.intersperse [[46]]).flatten

private def mapped (a : AddressFields) : Bool :=
  a.detail == .ipv6 && a.hi == 0 && a.lo.toNat / 2^32 == 65535

private def ipv4MappedChunks (a : AddressFields) : List (List UInt8) :=
  [textBytes "::ffff:"] ++ ipv4Chunks a ++ zone a

private def ipv6Words (a : AddressFields) : List UInt16 :=
  (List.range 8).map (fun i => addrV6Word a (UInt8.ofNat i))

private def runOfZeros (words : List UInt16) : Nat :=
  (words.takeWhile (· == 0)).length

private def longestZeros (words : List UInt16) : Nat × Nat :=
  (List.range words.length).foldl (fun (start, count) i =>
    let n := runOfZeros (words.drop i)
    if n ≥ 2 && n > count then (i, n) else (start, count)) (0, 0)

private def ipv6Chunks (a : AddressFields) : List (List UInt8) :=
  let words := ipv6Words a
  let piece := fun words : List UInt16 =>
    ((words.map fun w => (hexDigits w).map (fun b => [b])).intersperse [[58]]).flatten
  let (start, count) := longestZeros words
  let core := if count = 0 then piece words
    else piece (words.take start) ++ [[58,58]] ++ piece (words.drop (start + count))
  core ++ zone a

def addrAppendTo4 (a : AddressFields) (b : BinaryBuffer) : BinaryAppendResult :=
  appendTextChunks b (ipv4Chunks a)

def addrAppendTo4In6 (a : AddressFields) (b : BinaryBuffer) : BinaryAppendResult :=
  appendTextChunks b (ipv4MappedChunks a)

def addrAppendTo6 (a : AddressFields) (b : BinaryBuffer) : BinaryAppendResult :=
  appendTextChunks b (ipv6Chunks a)

def addrTextChunks (a : AddressFields) : List (List UInt8) :=
  match a.detail with
  | .invalid => []
  | .ipv4 => ipv4Chunks a
  | .ipv6 => if mapped a then ipv4MappedChunks a else ipv6Chunks a

def addrAppendTo (a : AddressFields) (b : BinaryBuffer) : BinaryAppendResult :=
  appendTextChunks b (addrTextChunks a)

def addrAppendText := addrAppendTo

def addrMarshalText (a : AddressFields) : List UInt8 := (addrTextChunks a).flatten

def addrString (a : AddressFields) : List UInt8 :=
  if a.detail = .invalid then textBytes "invalid IP" else addrMarshalText a

def addrStringExpanded (a : AddressFields) : List UInt8 :=
  if a.detail = .ipv6 then
    (((ipv6Words a).map hexPadDigits).intersperse [58]).flatten ++ (zone a).flatten
  else addrString a

private def invalid : AddressFields := ⟨.invalid, 0, 0, []⟩

theorem addrAppendTo_invalid (b : BinaryBuffer) :
    addrAppendTo invalid b = ⟨b.cells.take b.length, b.cells, true,
      b.nilSlice && (b.cells.take b.length).isEmpty⟩ := by
  simp [addrAppendTo, addrTextChunks, invalid, appendTextChunks]

theorem addrAppendText_same (a : AddressFields) (b : BinaryBuffer) :
    addrAppendText a b = addrAppendTo a b := rfl

theorem addrMarshalText_invalid : addrMarshalText invalid = [] := rfl
theorem addrString_invalid : addrString invalid = textBytes "invalid IP" := rfl
theorem addrStringExpanded_invalid : addrStringExpanded invalid = textBytes "invalid IP" := rfl

theorem addrAppendTo4_prefix (a : AddressFields) (b : BinaryBuffer) :
    (addrAppendTo4 a b).output = b.cells.take b.length ++ (ipv4Chunks a).flatten := rfl
theorem addrAppendTo4In6_prefix (a : AddressFields) (b : BinaryBuffer) :
    (addrAppendTo4In6 a b).output = b.cells.take b.length ++ (ipv4MappedChunks a).flatten := rfl
theorem addrAppendTo6_prefix (a : AddressFields) (b : BinaryBuffer) :
    (addrAppendTo6 a b).output = b.cells.take b.length ++ (ipv6Chunks a).flatten := rfl

end GoLean.Net.Netip
