import GoLean.Net.Netip.AddrFrom16

namespace GoLean.Net.Netip

/-- src/net/netip/netip.go:82 sets only byte 15 in the zero-filled array. -/
def ipv6Loopback : Addr := addrFrom16 fun i => if i = 15 then 1 else 0

theorem ipv6Loopback_spec : ipv6Loopback = ⟨0, 1, true, ""⟩ := by
  rfl


end GoLean.Net.Netip
