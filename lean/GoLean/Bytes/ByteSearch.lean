import GoLean.Slices.Index

namespace GoLean.Bytes

/-- The generic forward byte scan. None represents Go's -1.
    https://github.com/golang/go/blob/go1.27.0/src/bytes/bytes.go#L96-L100
    https://github.com/golang/go/blob/go1.27.0/src/internal/bytealg/indexbyte_generic.go#L13-L20 -/
def indexByte (xs : List UInt8) (key : UInt8) : Option Nat :=
  Slices.sliceIndex (· == ·) xs key

/-- Go scans from the back. Convert a match's reverse offset to its forward index.
    https://github.com/golang/go/blob/go1.27.0/src/bytes/bytes.go#L120-L124
    https://github.com/golang/go/blob/go1.27.0/src/internal/bytealg/lastindexbyte_generic.go#L7-L14 -/
def lastIndexByte (xs : List UInt8) (key : UInt8) : Option Nat :=
  (indexByte xs.reverse key).map (fun i => xs.length - 1 - i)

theorem indexByte_found (xs : List UInt8) (key : UInt8) (i : Nat) :
    indexByte xs key = some i ↔
      ∃ hi : i < xs.length, key = xs[i] ∧
        ∀ j (hj : j < i), key ≠ xs[j]'(Nat.lt_trans hj hi) := by
  simp only [indexByte, Slices.sliceIndex_found, beq_iff_eq]

theorem indexByte_missing (xs : List UInt8) (key : UInt8) :
    indexByte xs key = none ↔ key ∉ xs := by
  simp only [indexByte, Slices.sliceIndex_missing, beq_eq_false_iff_ne]
  constructor
  · intro absent mem
    exact absent key mem rfl
  · intro absent x mem eq
    exact absent (eq ▸ mem)

theorem lastIndexByte_found (xs : List UInt8) (key : UInt8) (i : Nat) :
    lastIndexByte xs key = some i ↔
      ∃ hi : i < xs.length, key = xs[i] ∧ ∀ j (hj : j < xs.length), i < j → key ≠ xs[j] := by
  rw [lastIndexByte, Option.map_eq_some_iff]
  constructor
  · rintro ⟨offset, hoffset, index⟩
    obtain ⟨hb, found, least⟩ := (indexByte_found xs.reverse key offset).mp hoffset
    have hlen : offset < xs.length := by simpa using hb
    refine ⟨by omega, ?_, ?_⟩
    · simpa only [List.getElem_reverse, index] using found
    · intro j hj larger
      have earlier : xs.length - 1 - j < offset := by omega
      have notFound := least _ earlier
      rw [List.getElem_reverse] at notFound
      have twice : xs.length - 1 - (xs.length - 1 - j) = j := by omega
      simpa only [twice] using notFound
  · rintro ⟨hi, found, greatest⟩
    let offset := xs.length - 1 - i
    have hb : offset < xs.reverse.length := by simp [offset]; omega
    refine ⟨offset, (indexByte_found xs.reverse key offset).mpr ⟨hb, ?_, ?_⟩, by omega⟩
    · rw [List.getElem_reverse]
      have twice : xs.length - 1 - offset = i := by omega
      simpa only [twice] using found
    · intro j earlier
      rw [List.getElem_reverse]
      exact greatest _ (by omega) (by omega)

theorem lastIndexByte_missing (xs : List UInt8) (key : UInt8) :
    lastIndexByte xs key = none ↔ key ∉ xs := by
  simp [lastIndexByte, indexByte_missing]

end GoLean.Bytes
