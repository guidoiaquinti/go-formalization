import GoLean.Net.Netip.AddrFrom16

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:77-79 fills three bytes before calling AddrFrom16. -/
def ipv6LinkLocalAllRouters : Addr :=
  addrFrom16 fun i =>
    if i = 0 then 0xff else if i = 1 || i = 15 then 2 else 0

theorem ipv6LinkLocalAllRouters_spec :
    ipv6LinkLocalAllRouters = ⟨0xff02000000000000, 2, true, ""⟩ := by
  rfl

end GoLean.Net.Netip
