import GoLean.Net.Netip.AddrUnmap
import GoLean.Net.Netip.AddrIs4

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:761-773 and uint128.go:52-55. -/
def addrPrev (ip : ZonedAddr) : ZonedAddr :=
  let atStart := if addrIs4 ip.detail then ip.lo.toUInt32 = 0 else ip.hi = 0 ∧ ip.lo = 0
  if atStart then ⟨.invalid, 0, 0, ""⟩
  else { ip with hi := ip.hi - (if ip.lo = 0 then 1 else 0), lo := ip.lo - 1 }

theorem addrPrev_spec (ip : ZonedAddr) :
    let atStart := (ip.detail = .ipv4 ∧ ip.lo.toNat % 2^32 = 0) ∨
      (ip.detail ≠ .ipv4 ∧ ip.hi = 0 ∧ ip.lo = 0)
    let result := addrPrev ip
    result.lo.toNat = (if atStart then 0 else (2^64 - 1 + ip.lo.toNat) % 2^64) ∧
    result.hi.toNat =
      (if atStart then 0 else (2^64 - (if ip.lo = 0 then 1 else 0) + ip.hi.toNat) % 2^64) ∧
    result.detail = (if atStart then .invalid else ip.detail) ∧
    result.zone = (if atStart then "" else ip.zone) := by
  have v4Start : ip.lo.toUInt32 = 0 ↔ ip.lo.toNat % 2^32 = 0 := by
    rw [← UInt32.toNat_inj]
    simp
  have start :
      (if addrIs4 ip.detail then ip.lo.toUInt32 = 0 else ip.hi = 0 ∧ ip.lo = 0) ↔
      (ip.detail = .ipv4 ∧ ip.lo.toNat % 2^32 = 0) ∨
      (ip.detail ≠ .ipv4 ∧ ip.hi = 0 ∧ ip.lo = 0) := by
    cases ip.detail <;> simp [addrIs4, v4Start]
  simp only [addrPrev, start]
  split
  · simp
  · by_cases h : ip.lo = 0 <;> simp [h, UInt64.toNat_sub]

end GoLean.Net.Netip
