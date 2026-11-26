import GoLean.Net.Netip.AddrUnmap

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:1058-1065 stores the address and port separately. -/
structure AddrPort where
  ip : ZonedAddr
  port : UInt16

/-- Go 1.27.0 src/net/netip/netip.go:1071 reads the stored port. -/
def addrPortPort (p : AddrPort) : UInt16 := p.port

theorem addrPortPort_spec (ip : ZonedAddr) (port : UInt16) :
    addrPortPort ⟨ip, port⟩ = port := rfl


end GoLean.Net.Netip
