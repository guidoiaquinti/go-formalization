import GoLean.Net.Netip.AddrPortFrom

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:1067-1068 returns the stored address. -/
def addrPortAddr (p : AddrPort) : ZonedAddr := p.ip

theorem addrPortAddr_spec (ip : ZonedAddr) (port : UInt16) :
    addrPortAddr (addrPortFrom ip port) = ip := rfl


end GoLean.Net.Netip
