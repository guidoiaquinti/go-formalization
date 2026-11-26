import GoLean.Net.Netip.AddrFrom16

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:73-75 fills three bytes before calling AddrFrom16. -/
def ipv6LinkLocalAllNodes : Addr :=
  addrFrom16 fun i =>
    if i = 0 then 0xff else if i = 1 then 2 else if i = 15 then 1 else 0

theorem ipv6LinkLocalAllNodes_spec :
    ipv6LinkLocalAllNodes = ⟨0xff02000000000000, 1, true, ""⟩ := rfl

end GoLean.Net.Netip
