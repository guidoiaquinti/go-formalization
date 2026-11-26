import GoLean.Net.Netip.AddrUnmap
import GoLean.Net.Netip.AddrAs16
import GoLean.Net.Netip.AddrIs4

namespace GoLean.Net.Netip

private def privateV4 (a b : UInt8) : Bool :=
  a == 10 || (a == 172 && (b &&& 0xf0) == 16) || (a == 192 && b == 168)

/-- Go 1.27.0 src/net/netip/netip.go:643-664. The byte accessors at
360-370 read the IPv4 payload and the high IPv6 octet, respectively. -/
def addrIsPrivate (input : ZonedAddr) : Bool :=
  let ip := if addrIs4In6 input.detail input.hi input.lo then addrUnmap input else input
  if addrIs4 ip.detail then privateV4 (addrAs16 ip.hi ip.lo 12) (addrAs16 ip.hi ip.lo 13)
  else if addrIs6 ip.detail then (addrAs16 ip.hi ip.lo 0 &&& 0xfe) == 0xfc
  else false

theorem addrIsPrivate_spec (ip : ZonedAddr) :
    addrIsPrivate ip = true ↔
      ((ip.detail = .ipv4 ∨
          (ip.detail = .ipv6 ∧ ip.hi = 0 ∧ ip.lo.toNat / 2^32 = 0xffff)) ∧
        privateV4 (addrAs16 ip.hi ip.lo 12) (addrAs16 ip.hi ip.lo 13) = true) ∨
      (ip.detail = .ipv6 ∧ (addrAs16 ip.hi ip.lo 0 &&& 0xfe) = 0xfc) := by
  have zeroOctet : addrAs16 0 ip.lo 0 = 0 := by simp [addrAs16]
  cases h : ip.detail <;>
    simp [addrIsPrivate, addrUnmap, addrIs4In6_spec, addrIs4, addrIs6, h]
  by_cases high : ip.hi = 0 <;> by_cases low : ip.lo.toNat / 2^32 = 0xffff <;>
    simp_all

end GoLean.Net.Netip
