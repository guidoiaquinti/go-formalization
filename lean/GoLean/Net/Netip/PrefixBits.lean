import GoLean.Net.Netip.PrefixAddr

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:1307-1310 widens the tag before subtracting. -/
def prefixBits (p : Prefix) : Int := Int.ofNat p.bitsPlusOne.toNat - 1

theorem prefixBits_spec (p : Prefix) (bits : Int) :
    prefixBits p = bits ↔ Int.ofNat p.bitsPlusOne.toNat = bits + 1 := by
  unfold prefixBits
  omega


end GoLean.Net.Netip
