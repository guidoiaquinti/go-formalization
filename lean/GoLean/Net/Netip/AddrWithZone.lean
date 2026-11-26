import GoLean.Net.Netip.AddrUnmap

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:488-501.
Both handle assignments are represented by their observable zone contents. -/
def addrWithZone (ip : ZonedAddr) (zone : String) : ZonedAddr :=
  if !addrIs6 ip.detail then ip
  else if zone = "" then { ip with zone := "" }
  else { ip with zone := zone }

theorem addrWithZone_spec (ip : ZonedAddr) (zone : String) :
    addrWithZone ip zone =
      if ip.detail = .ipv6 then { ip with zone := zone } else ip := by
  cases h : ip.detail <;> simp [addrWithZone, addrIs6, h]

end GoLean.Net.Netip
