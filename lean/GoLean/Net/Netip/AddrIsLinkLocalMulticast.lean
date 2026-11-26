import GoLean.Net.Netip.AddrUnmap
import GoLean.Net.Netip.AddrAs16
import GoLean.Net.Netip.AddrIs4

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:587-603.
The v4 helper selects octets 12 through 14. The v6u16 helper selects the high 16 bits. -/
def addrIsLinkLocalMulticast (input : ZonedAddr) : Bool :=
  let ip := if addrIs4In6 input.detail input.hi input.lo then addrUnmap input else input
  if addrIs4 ip.detail then addrAs16 ip.hi ip.lo 12 == 224 && addrAs16 ip.hi ip.lo 13 == 0 &&
    addrAs16 ip.hi ip.lo 14 == 0
  else if addrIs6 ip.detail then ((ip.hi >>> 48).toUInt16 &&& 0xff0f) == 0xff02
  else false

theorem addrIsLinkLocalMulticast_spec (ip : ZonedAddr) :
    addrIsLinkLocalMulticast ip = true ↔
      ((ip.detail = .ipv4 ∨
          (ip.detail = .ipv6 ∧ ip.hi = 0 ∧ ip.lo.toNat / 2^32 = 0xffff)) ∧
        addrAs16 ip.hi ip.lo 12 = 224 ∧
        addrAs16 ip.hi ip.lo 13 = 0 ∧ addrAs16 ip.hi ip.lo 14 = 0) ∨
      (ip.detail = .ipv6 ∧
        ((ip.hi >>> 48).toUInt16 &&& 0xff0f) = 0xff02) := by
  cases h : ip.detail <;>
    simp [addrIsLinkLocalMulticast, addrUnmap, addrIs4In6_spec, addrIs4, addrIs6, h, and_assoc]
  by_cases high : ip.hi = 0 <;> by_cases low : ip.lo.toNat / 2^32 = 0xffff <;>
    simp_all

end GoLean.Net.Netip
