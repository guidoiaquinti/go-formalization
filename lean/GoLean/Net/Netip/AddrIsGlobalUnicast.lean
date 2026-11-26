import GoLean.Net.Netip.AddrIsValid
import GoLean.Net.Netip.AddrIsUnspecified
import GoLean.Net.Netip.AddrIsLoopback
import GoLean.Net.Netip.AddrIsMulticast
import GoLean.Net.Netip.AddrIsLinkLocalUnicast

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:614-634. -/
def addrIsGlobalUnicast (input : ZonedAddr) : Bool :=
  if !addrIsValid input.detail then false
  else
    let ip := if addrIs4In6 input.detail input.hi input.lo then addrUnmap input else input
    if addrIs4 ip.detail &&
      (matchesAddr ip ipv4Unspecified || matchesAddr ip (addrFrom4 fun _ => 255)) then false
    else !matchesAddr ip ipv6Unspecified &&
      !addrIsLoopback ip && !addrIsMulticast ip && !addrIsLinkLocalUnicast ip

theorem addrIsGlobalUnicast_spec (input : ZonedAddr) :
    addrIsGlobalUnicast input = true ↔
      let ip := addrUnmap input
      input.detail ≠ .invalid ∧
        (ip.detail = .ipv4 → matchesAddr ip (addrFrom4 fun _ => 255) = false) ∧
        addrIsUnspecified ip = false ∧ addrIsLoopback ip = false ∧
        addrIsMulticast ip = false ∧ addrIsLinkLocalUnicast ip = false := by
  have normalized :
      (if addrIs4In6 input.detail input.hi input.lo then addrUnmap input else input) =
        addrUnmap input := by
    unfold addrUnmap
    split <;> simp_all
  unfold addrIsGlobalUnicast
  rw [normalized]
  have valid := addrIsValid_spec input.detail
  generalize addrUnmap input = ip
  cases hv : addrIsValid input.detail <;>
    cases hd : ip.detail <;>
    simp_all [addrIsValid_spec, addrIs4, addrIsUnspecified, matchesAddr,
      ipv4Unspecified_spec, ipv6Unspecified_spec, and_assoc] <;> grind

end GoLean.Net.Netip
