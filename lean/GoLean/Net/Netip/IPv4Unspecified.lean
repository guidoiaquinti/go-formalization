import GoLean.Net.Netip.AddrFrom4

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:87-88 passes a zero-filled array. -/
def ipv4Unspecified : Addr := addrFrom4 fun _ => 0

theorem ipv4Unspecified_spec :
    ipv4Unspecified = ⟨0, 0xffff00000000, false, ""⟩ := rfl

end GoLean.Net.Netip
