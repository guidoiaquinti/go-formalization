import Init.Data.BitVec.Lemmas
import Init.Data.UInt.Lemmas

namespace GoLean.Net.Netip

/-- Observable contents of the two address words and the interned address detail
in Go 1.27.0 src/net/netip/netip.go:37-70. -/
structure Addr where
  hi : UInt64
  lo : UInt64
  isV6 : Bool
  zone : String

/-- src/internal/byteorder/byteorder.go:120-123, on an eight-byte window.
The fixed input size guarantees that every byte access succeeds. -/
private def beUint64 (b : Fin 8 → UInt8) : UInt64 :=
  (b 7).toUInt64 ||| ((b 6).toUInt64 <<< 8) ||| ((b 5).toUInt64 <<< 16) |||
  ((b 4).toUInt64 <<< 24) ||| ((b 3).toUInt64 <<< 32) ||| ((b 2).toUInt64 <<< 40) |||
  ((b 1).toUInt64 <<< 48) ||| ((b 0).toUInt64 <<< 56)

/-- src/net/netip/netip.go:101-109. Decode the two halves and attach z6noz. -/
def addrFrom16 (b : Fin 16 → UInt8) : Addr :=
  ⟨beUint64 (fun i => b ⟨i.val, by omega⟩),
   beUint64 (fun i => b ⟨8 + i.val, by omega⟩), true, ""⟩

/-- Concatenate the eight octets in network order. -/
def networkWord (b : Fin 8 → UInt8) : UInt64 :=
  UInt64.ofBitVec ((b 0).toBitVec ++ (b 1).toBitVec ++ (b 2).toBitVec ++
    (b 3).toBitVec ++ (b 4).toBitVec ++ (b 5).toBitVec ++ (b 6).toBitVec ++ (b 7).toBitVec)

/-- Split a finite set of positions while leaving the input word symbolic. -/
theorem positions_succ {p : Nat → Prop} (first : p 0)
    (rest : ∀ i < n, p (i + 1)) : ∀ i < n + 1, p i := by
  intro i hi
  cases i with
  | zero => exact first
  | succ i => exact rest i (by omega)

private theorem beUint64_spec (b : Fin 8 → UInt8) : beUint64 b = networkWord b := by
  apply UInt64.eq_of_toBitVec_eq
  change (((b 7).toBitVec.setWidth 64) ||| ((b 6).toBitVec.setWidth 64 <<< 8) |||
    ((b 5).toBitVec.setWidth 64 <<< 16) ||| ((b 4).toBitVec.setWidth 64 <<< 24) |||
    ((b 3).toBitVec.setWidth 64 <<< 32) ||| ((b 2).toBitVec.setWidth 64 <<< 40) |||
    ((b 1).toBitVec.setWidth 64 <<< 48) ||| ((b 0).toBitVec.setWidth 64 <<< 56)) =
    ((b 0).toBitVec ++ (b 1).toBitVec ++ (b 2).toBitVec ++ (b 3).toBitVec ++
      (b 4).toBitVec ++ (b 5).toBitVec ++ (b 6).toBitVec ++ (b 7).toBitVec)
  apply BitVec.eq_of_getLsbD_eq
  simp only [BitVec.getLsbD_append]
  iterate 64
    apply positions_succ
    · simp +decide
  omega

theorem addrFrom16_spec (b : Fin 16 → UInt8) :
    addrFrom16 b = ⟨networkWord (fun i => b ⟨i.val, by omega⟩),
      networkWord (fun i => b ⟨8 + i.val, by omega⟩), true, ""⟩ := by
  simp only [addrFrom16, beUint64_spec]

end GoLean.Net.Netip
