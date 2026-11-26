import GoLean.Net.Netip.AddrUnmap
import GoLean.Net.Netip.AddrIs4

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:745-758 and uint128.go:58-61.
The low word wraps before carry and family-specific overflow are tested. -/
def addrNext (ip : ZonedAddr) : ZonedAddr :=
  let lo := ip.lo + 1
  let hi := ip.hi + (if lo = 0 then 1 else 0)
  let overflow := if addrIs4 ip.detail then lo.toUInt32 = 0 else hi = 0 ∧ lo = 0
  if overflow then ⟨.invalid, 0, 0, ""⟩ else { ip with hi := hi, lo := lo }

theorem addrNext_spec (ip : ZonedAddr) :
    let wrapped := (ip.detail = .ipv4 ∧ ip.lo.toNat % 2^32 = 2^32 - 1) ∨
      (ip.detail ≠ .ipv4 ∧ ip.hi.toNat = 2^64 - 1 ∧ ip.lo.toNat = 2^64 - 1)
    (addrNext ip).lo.toNat = (if wrapped then 0 else (ip.lo.toNat + 1) % 2^64) ∧
    (addrNext ip).hi.toNat =
      (if wrapped then 0 else (ip.hi.toNat + (ip.lo.toNat + 1) / 2^64) % 2^64) ∧
    (addrNext ip).detail = (if wrapped then .invalid else ip.detail) ∧
    (addrNext ip).zone = (if wrapped then "" else ip.zone) := by
  have loBound := ip.lo.toNat_lt
  have hiBound := ip.hi.toNat_lt
  have carry : ip.lo + 1 = 0 ↔ ip.lo.toNat = 2^64 - 1 := by
    rw [← UInt64.toNat_inj]
    simp [UInt64.toNat_add]
    omega
  have quotient : (ip.lo.toNat + 1) / 2^64 =
      (if ip.lo + 1 = 0 then 1 else 0) := by
    split <;> rename_i h
    · have := carry.mp h
      omega
    · have : ip.lo.toNat ≠ 2^64 - 1 := fun e => h (carry.mpr e)
      omega
  have v4Wrap : (ip.lo + 1).toUInt32 = 0 ↔ ip.lo.toNat % 2^32 = 2^32 - 1 := by
    rw [← UInt32.toNat_inj]
    simp
    omega
  have v6Wrap : (ip.hi + (if ip.lo + 1 = 0 then 1 else 0) = 0 ∧ ip.lo + 1 = 0) ↔
      ip.hi.toNat = 2^64 - 1 ∧ ip.lo.toNat = 2^64 - 1 := by
    by_cases h : ip.lo + 1 = 0
    · simp only [h, ite_true, and_true, (carry.mp h)]
      rw [← UInt64.toNat_inj]
      simp [UInt64.toNat_add]
      omega
    · simp [h, show ip.lo.toNat ≠ 2^64 - 1 from fun e => h (carry.mpr e)]
  have v4WrapReduced : ip.lo.toUInt32 + 1 = 0 ↔ ip.lo.toNat % 2^32 = 2^32 - 1 := by
    simpa using v4Wrap
  have overflow :
      (if addrIs4 ip.detail then (ip.lo + 1).toUInt32 = 0
        else ip.hi + (if ip.lo + 1 = 0 then 1 else 0) = 0 ∧ ip.lo + 1 = 0) ↔
      (ip.detail = .ipv4 ∧ ip.lo.toNat % 2^32 = 2^32 - 1) ∨
      (ip.detail ≠ .ipv4 ∧ ip.hi.toNat = 2^64 - 1 ∧ ip.lo.toNat = 2^64 - 1) := by
    cases h : ip.detail <;> simp [addrIs4, v4WrapReduced, v6Wrap]
  simp only [addrNext, overflow]
  split
  · simp
  · by_cases h : ip.lo + 1 = 0 <;> simp [h, quotient, UInt64.toNat_add]
    simp [carry.mp h]

end GoLean.Net.Netip
