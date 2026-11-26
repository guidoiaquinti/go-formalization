import GoLean.Net.Netip.AddrUnmap

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:1277-1283. -/
structure Prefix where
  ip : ZonedAddr
  bitsPlusOne : UInt8

/-- Go 1.27.0 src/net/netip/netip.go:1304-1305. -/
def prefixAddr (p : Prefix) : ZonedAddr := p.ip

theorem prefixAddr_spec (ip : ZonedAddr) (bitsPlusOne : UInt8) :
    prefixAddr ⟨ip, bitsPlusOne⟩ = ip := rfl


end GoLean.Net.Netip
