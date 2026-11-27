import GoLean.Net.Netip.BinaryCodec

namespace GoLean.Net.Netip

private def digit (x : Nat) : UInt8 :=
  UInt8.ofNat (if x < 10 then 48 + x else 87 + x)

/-- Go 1.27.0 src/net/netip/netip.go:838. Each digit is appended separately. -/
def decimalDigits (x : UInt8) : List UInt8 :=
  let n := x.toNat
  (if n ≥ 100 then [digit (n / 100)] else []) ++
  (if n ≥ 10 then [digit (n / 10 % 10)] else []) ++ [digit (n % 10)]

/-- Go 1.27.0 src/net/netip/netip.go:852. Emit lowercase hex without leading zeroes. -/
def hexDigits (x : UInt16) : List UInt8 :=
  let n := x.toNat
  (if n ≥ 4096 then [digit (n / 4096)] else []) ++
  (if n ≥ 256 then [digit (n / 256 % 16)] else []) ++
  (if n ≥ 16 then [digit (n / 16 % 16)] else []) ++ [digit (n % 16)]

def hexPadDigits (x : UInt16) : List UInt8 :=
  let n := x.toNat
  [digit (n / 4096), digit (n / 256 % 16), digit (n / 16 % 16), digit (n % 16)]

private def emitDigits (b : BinaryBuffer) (written : List UInt8) (atomic : Bool) :
    BinaryAppendResult :=
  let limit := b.cells.length - b.length
  let reused := written.length ≤ limit
  let retained := if atomic && !reused then [] else written.take limit
  ⟨b.cells.take b.length ++ written,
   b.cells.take b.length ++ retained ++ b.cells.drop (b.length + retained.length),
   reused, false⟩

def appendDecimal (b : BinaryBuffer) (x : UInt8) : BinaryAppendResult :=
  emitDigits b (decimalDigits x) false

def appendHex (b : BinaryBuffer) (x : UInt16) : BinaryAppendResult :=
  emitDigits b (hexDigits x) false

/-- Go 1.27.0 src/net/netip/netip.go:869. A single four-byte append. -/
def appendHexPad (b : BinaryBuffer) (x : UInt16) : BinaryAppendResult :=
  emitDigits b (hexPadDigits x) true

theorem appendDecimal_zero (b : BinaryBuffer) :
    (appendDecimal b 0).output = b.cells.take b.length ++ [48] := rfl

theorem appendHex_zero (b : BinaryBuffer) :
    (appendHex b 0).output = b.cells.take b.length ++ [48] := rfl

theorem appendHexPad_width (b : BinaryBuffer) (x : UInt16) :
    (appendHexPad b x).output.length = (b.cells.take b.length).length + 4 := by
  simp [appendHexPad, emitDigits, hexPadDigits]

end GoLean.Net.Netip
