import GoLean.Net.Netip.AddrUnmap

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:407-413. The zero handle has no detail to read. -/
def addrZone (ip : ZonedAddr) : String :=
  if ip.detail = .invalid then "" else ip.zone

theorem addrZone_spec (ip : ZonedAddr) (ipv4Unscoped : ip.detail = .ipv4 → ip.zone = "") :
    addrZone ip = (if ip.detail = .ipv6 then ip.zone else "") := by
  cases h : ip.detail <;> simp_all [addrZone]

end GoLean.Net.Netip
