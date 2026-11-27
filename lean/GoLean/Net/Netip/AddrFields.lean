import GoLean.Net.Netip.AddrFrom4
import GoLean.Net.Netip.AddrBitLen

namespace GoLean.Net.Netip

structure AddressFields where
  detail : AddrDetailClass
  hi : UInt64
  lo : UInt64
  zone : List UInt8

/-- Go 1.27.0 src/net/netip/netip.go:348. Exact-length slice conversions
choose the existing fixed-width address constructors. -/
def addrFromSlice (raw : List UInt8) : Option Addr :=
  if raw.length = 4 then some (addrFrom4 fun i => raw[i.val]!)
  else if raw.length = 16 then some (addrFrom16 fun i => raw[i.val]!) else none

/-- The shift count in src/net/netip/netip.go:360 uses uint8 arithmetic.
A Go shift of at least the operand width yields zero. -/
def addrV4Byte (ip : AddressFields) (i : UInt8) : UInt8 :=
  let shift := ((3 - i) * 8).toNat
  (ip.lo.toBitVec >>> shift).toNat.toUInt8

def addrV6Byte (ip : AddressFields) (i : UInt8) : UInt8 :=
  let index := i.toNat
  let word := if index / 8 % 2 = 0 then ip.hi else ip.lo
  (word.toBitVec >>> ((7 - index % 8) * 8)).toNat.toUInt8

def addrV6Word (ip : AddressFields) (i : UInt8) : UInt16 :=
  let index := i.toNat
  let word := if index / 4 % 2 = 0 then ip.hi else ip.lo
  (word.toBitVec >>> ((3 - index % 4) * 16)).toNat.toUInt16

def addrIsZero (ip : AddressFields) : Bool := ip.detail == .invalid

def addrHasZone (ip : AddressFields) : Bool :=
  ip.detail == .ipv6 && !ip.zone.isEmpty

def addrWithoutZone (ip : AddressFields) : AddressFields :=
  if ip.detail = .ipv6 then { ip with zone := [] } else ip

theorem addrFromSlice_invalid (raw : List UInt8) (h4 : raw.length ≠ 4) (h16 : raw.length ≠ 16) :
    addrFromSlice raw = none := by simp [addrFromSlice, h4, h16]

theorem addrV4Byte_wrap (ip : AddressFields) : addrV4Byte ip 32 = addrV4Byte ip 0 := rfl
theorem addrV6Byte_wrap (ip : AddressFields) : addrV6Byte ip 16 = addrV6Byte ip 0 := rfl
theorem addrV6Word_wrap (ip : AddressFields) : addrV6Word ip 8 = addrV6Word ip 0 := rfl
theorem addrIsZero_ipv4 (hi lo : UInt64) : addrIsZero ⟨.ipv4, hi, lo, []⟩ = false := rfl

theorem addrHasZone_empty (detail : AddrDetailClass) (hi lo : UInt64) :
    addrHasZone ⟨detail, hi, lo, []⟩ = false := by simp [addrHasZone]

theorem addrWithoutZone_ipv6 (hi lo : UInt64) (zone : List UInt8) :
    addrWithoutZone ⟨.ipv6, hi, lo, zone⟩ = ⟨.ipv6, hi, lo, []⟩ := rfl

end GoLean.Net.Netip
