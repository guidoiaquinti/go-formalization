import GoLean.Net.Netip.AddrUnmap

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:575-584 and v6u16 at 372-374. -/
def addrIsInterfaceLocalMulticast (ip : ZonedAddr) : Bool :=
  if addrIs6 ip.detail && !addrIs4In6 ip.detail ip.hi ip.lo then
    ((ip.hi >>> 48).toUInt16 &&& 0xff0f) == 0xff01
  else false

theorem addrIsInterfaceLocalMulticast_spec (ip : ZonedAddr) :
    addrIsInterfaceLocalMulticast ip = true ↔
      ip.detail = .ipv6 ∧ ((ip.hi >>> 48).toUInt16 &&& 0xff0f) = 0xff01 := by
  cases h : ip.detail <;> simp [addrIsInterfaceLocalMulticast, addrIs6, h]
  by_cases high : ip.hi = 0 <;> simp_all [addrIs4In6]


end GoLean.Net.Netip
