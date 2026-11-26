import GoLean.Net.Netip.AddrIs4In6

namespace GoLean.Net.Netip

structure ZonedAddr where
  detail : AddrDetailClass
  hi : UInt64
  lo : UInt64
  zone : String

/-- Go 1.27.0 src/net/netip/netip.go:477-486 replaces only the detail handle. -/
def addrUnmap (ip : ZonedAddr) : ZonedAddr :=
  if addrIs4In6 ip.detail ip.hi ip.lo then
    { ip with detail := .ipv4, zone := "" }
  else ip

theorem addrUnmap_spec (ip : ZonedAddr) :
    addrUnmap ip =
      if ip.detail = .ipv6 ∧ ip.hi = 0 ∧ ip.lo.toNat / 2^32 = 0xffff then
        ⟨.ipv4, ip.hi, ip.lo, ""⟩
      else ip := by
  simp only [addrUnmap, addrIs4In6_spec]

end GoLean.Net.Netip
