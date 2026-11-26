import GoLean.Net.Netip.AddrAs16
import GoLean.Net.Netip.AddrIs4In6

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:712-724.
BEPutUint32 of the truncated low word is the final four big-endian bytes. -/
def addrAs4 (detail : AddrDetailClass) (hi lo : UInt64) : Except String (Fin 4 → UInt8) :=
  if (detail == .ipv4) || addrIs4In6 detail hi lo then
    .ok fun i => addrAs16 hi lo ⟨12 + i.val, by omega⟩
  else if detail = .invalid then .error "As4 called on IP zero value"
  else .error "As4 called on IPv6 address"

private theorem tailByte (hi lo : UInt64) (i : Fin 4) :
    addrAs16 hi lo ⟨12 + i.val, by omega⟩ =
      UInt8.ofNat (lo.toNat / 256 ^ (3 - i.val) % 256) := by
  apply UInt8.toNat_inj.1
  rw [addrAs16_spec]
  have offset : (12 + i.val) % 8 = 4 + i.val := by omega
  have exponent : 7 - (4 + i.val) = 3 - i.val := by omega
  simp [offset, exponent, show ¬12 + i.val < 8 by omega]

theorem addrAs4_spec (detail : AddrDetailClass) (hi lo : UInt64) :
    addrAs4 detail hi lo =
      if detail = .ipv4 ∨ (detail = .ipv6 ∧ hi = 0 ∧ lo.toNat / 2^32 = 0xffff) then
        .ok (fun i => UInt8.ofNat (lo.toNat / 256 ^ (3 - i.val) % 256))
      else if detail = .invalid then .error "As4 called on IP zero value"
      else .error "As4 called on IPv6 address" := by
  simp [addrAs4, addrIs4In6_spec, tailByte]

end GoLean.Net.Netip
