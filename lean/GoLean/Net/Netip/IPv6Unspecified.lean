import GoLean.Net.Netip.AddrFrom16

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:65-71 and 84-85.
The address words take their zero values and z6noz carries the IPv6 flag and empty zone. -/
def ipv6Unspecified : Addr :=
  { hi := 0, lo := 0, isV6 := true, zone := "" }

theorem ipv6Unspecified_spec : ipv6Unspecified = ⟨0, 0, true, ""⟩ := rfl

end GoLean.Net.Netip
