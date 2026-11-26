import GoLean.Net.Netip.AddrFrom16

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:65-71 and 90-96.
Store the IPv4 bytes below the mapped prefix and attach the z4 detail. -/
def addrFrom4 (b : Fin 4 → UInt8) : Addr :=
  ⟨0, (0xffff00000000 : UInt64) ||| ((b 0).toUInt64 <<< 24) ||| ((b 1).toUInt64 <<< 16) |||
    ((b 2).toUInt64 <<< 8) ||| (b 3).toUInt64, false, ""⟩

private def mappedWord (b : Fin 4 → UInt8) : BitVec 64 :=
  0x0000ffff#32 ++ (b 0).toBitVec ++ (b 1).toBitVec ++ (b 2).toBitVec ++ (b 3).toBitVec

theorem addrFrom4_spec (b : Fin 4 → UInt8) :
    addrFrom4 b = ⟨0, UInt64.ofBitVec (mappedWord b), false, ""⟩ := by
  have low : (addrFrom4 b).lo = UInt64.ofBitVec (mappedWord b) := by
    apply UInt64.eq_of_toBitVec_eq
    change (0xffff00000000#64 ||| ((b 0).toBitVec.setWidth 64 <<< 24) |||
      ((b 1).toBitVec.setWidth 64 <<< 16) ||| ((b 2).toBitVec.setWidth 64 <<< 8) |||
      (b 3).toBitVec.setWidth 64) = mappedWord b
    unfold mappedWord
    apply BitVec.eq_of_getLsbD_eq
    simp only [BitVec.getLsbD_append]
    iterate 64
      apply positions_succ
      · simp +decide
    omega
  change Addr.mk 0 (addrFrom4 b).lo false "" = _
  rw [low]

end GoLean.Net.Netip
