import GoLean.Slices.Grow
import GoLean.Slices.Reverse

namespace GoLean.Slices

inductive InsertValues (α : Type) where
  | independent (values : List α)
  | alias (offset length : Nat)

def InsertValues.read (v : InsertValues α) (cells : List α) : List α :=
  match v with
  | .independent values => values
  | .alias offset length => (cells.drop offset).take length

private def InsertValues.overlaps (v : InsertValues α) (start stop : Nat) : Bool :=
  match v with
  | .independent _ => false
  | .alias offset length => start < offset + length && offset < stop

def writeCells (cells : List α) (start : Nat) (values : List α) : List α :=
  cells.take start ++ values ++ cells.drop (start + values.length)

private def rotateInsertTail (tail : List α) (oldLength : Nat) : List α :=
  sliceReverse (sliceReverse (tail.take oldLength) ++ sliceReverse (tail.drop oldLength))

/-- Models the original slice's capacity cells and the two in-place copy
    orders, including the overlap branch's three reversals. Result capacity
    beyond the new length after allocating is abstracted. Go 1.27.0
    src/slices/slices.go:136–212, 445–469. -/
def sliceInsert (s : GrowView α) (index : Int) (source : InsertValues α) :
    Except Unit (GrowObservation α) :=
  let i := index.toNat
  if index < 0 ∨ s.length < i then .error ()
  else
    let values := source.read s.cells
    let m := values.length
    let n := s.length
    let total := n + m
    let suffix := (s.cells.drop i).take (n - i)
    if m = 0 then .ok ⟨s, false⟩
    else if total > s.cells.length then
      .ok ⟨⟨s.cells.take i ++ values ++ suffix, total, false⟩, true⟩
    else if i = n then .ok ⟨⟨writeCells s.cells n values, total, false⟩, false⟩
    else if !source.overlaps (i + m) total then
      let shifted := writeCells s.cells (i + m) suffix
      .ok ⟨⟨writeCells shifted i (source.read shifted), total, false⟩, false⟩
    else
      let copied := writeCells s.cells n values
      let tail := (copied.drop i).take (total - i)
      .ok ⟨⟨writeCells copied i (rotateInsertTail tail (n - i)), total, false⟩, false⟩

theorem sliceInsert_bounds (s : GrowView α) (index : Int) (source : InsertValues α)
    (outOfBounds : index < 0 ∨ s.length < index.toNat) :
    sliceInsert s index source = .error () := by
  rw [sliceInsert, ite_eq_left outOfBounds]

theorem sliceInsert_empty (s : GrowView α) (index : Nat) (inBounds : index ≤ s.length) :
    sliceInsert s index (.independent []) = .ok ⟨s, false⟩ := by
  simp [sliceInsert, InsertValues.read, Nat.not_lt_of_ge inBounds]

theorem sliceInsert_allocate (s : GrowView α) (index : Nat) (source : InsertValues α)
    (inBounds : index ≤ s.length) (nonempty : (source.read s.cells).length ≠ 0)
    (insufficient : s.cells.length < s.length + (source.read s.cells).length) :
    sliceInsert s index source =
      .ok ⟨⟨s.cells.take index ++ source.read s.cells ++
        (s.cells.take s.length).drop index, s.length + (source.read s.cells).length, false⟩, true⟩ := by
  have suffix : (s.cells.drop index).take (s.length - index) = (s.cells.take s.length).drop index := by
    rw [List.take_drop, Nat.add_sub_cancel' inBounds]
  simp [sliceInsert, Nat.not_lt_of_ge inBounds, nonempty, insufficient, suffix]

private theorem rotateInsertTail_spec (suffix inserted : List α) :
    rotateInsertTail (suffix ++ inserted) suffix.length = inserted ++ suffix := by
  simp [rotateInsertTail, sliceReverse_spec]

theorem sliceInsert_overlap (before suffix spare : List α) (source : InsertValues α) (nilSlice : Bool)
    (hasSuffix : suffix ≠ [])
    (nonempty : (source.read (before ++ suffix ++ spare)).length ≠ 0)
    (fits : (source.read (before ++ suffix ++ spare)).length ≤ spare.length)
    (overlap : source.overlaps (before.length + (source.read (before ++ suffix ++ spare)).length)
      (before.length + suffix.length + (source.read (before ++ suffix ++ spare)).length) = true) :
    let values := source.read (before ++ suffix ++ spare)
    sliceInsert ⟨before ++ suffix ++ spare, before.length + suffix.length, nilSlice⟩
        before.length source =
      .ok ⟨⟨before ++ values ++ suffix ++ spare.drop values.length,
        before.length + suffix.length + values.length, false⟩, false⟩ := by
  dsimp only
  generalize measured : source.read (before ++ suffix ++ spare) = values at *
  have hlen : suffix.length ≠ 0 := by simpa using hasSuffix
  have cap : ¬ before.length + suffix.length + values.length > (before ++ suffix ++ spare).length := by
    simp only [List.length_append]
    omega
  have notEnd : before.length ≠ before.length + suffix.length := by omega
  have splitAtEnd : (before ++ suffix ++ spare).take (before.length + suffix.length) = before ++ suffix := by
    simpa only [List.length_append] using List.take_append_length (l₁ := before ++ suffix) (l₂ := spare)
  have spareAfter : (before ++ suffix ++ spare).drop
      (before.length + suffix.length + values.length) = spare.drop values.length := by
    rw [← List.drop_drop]
    have partition := List.drop_append_length (l₁ := before ++ suffix) (l₂ := spare)
    simp only [List.length_append] at partition
    rw [partition]
  have prefixKept : (before ++ suffix ++ values ++ spare.drop values.length).take before.length = before := by
    simp [List.append_assoc]
  have middle : ((before ++ suffix ++ values ++ spare.drop values.length).drop before.length).take
      (suffix.length + values.length) = suffix ++ values := by
    rw [show before ++ suffix ++ values ++ spare.drop values.length =
        before ++ ((suffix ++ values) ++ spare.drop values.length) by simp [List.append_assoc]]
    rw [List.drop_append_length]
    simpa only [List.length_append] using
      List.take_append_length (l₁ := suffix ++ values) (l₂ := spare.drop values.length)
  have rest : (before ++ suffix ++ values ++ spare.drop values.length).drop
      (before.length + (values ++ suffix).length) = spare.drop values.length := by
    rw [show before.length + (values ++ suffix).length = (before ++ suffix ++ values).length by
      simp only [List.length_append]; omega, List.drop_append_length]
  change sliceInsert ⟨before ++ suffix ++ spare, before.length + suffix.length, nilSlice⟩
      before.length source = _
  have lower : ¬ before.length + suffix.length < before.length := by omega
  have count : before.length + suffix.length + values.length - before.length =
      suffix.length + values.length := by omega
  have oldCount : before.length + suffix.length - before.length = suffix.length := by omega
  have nonnegative : ¬ (before.length : Int) < 0 := by omega
  simp only [sliceInsert, Int.toNat_natCast, nonnegative, lower, or_self, ite_false,
    measured, nonempty, cap, notEnd, overlap, Bool.not_true, Bool.false_eq_true,
    writeCells, splitAtEnd, spareAfter, count, oldCount, prefixKept, middle,
    rotateInsertTail_spec, rest]
  simp only [List.append_assoc]

end GoLean.Slices
