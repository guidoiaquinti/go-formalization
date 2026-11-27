import GoLean.Slices.Insert
import GoLean.Slices.Overlaps
import GoLean.Slices.RotateRight

namespace GoLean.Slices

private def overlapsTarget (source : InsertValues α) (start stop : Nat) : Bool :=
  match source with
  | .independent _ => false
  | .alias offset length => sliceOverlaps offset length start (stop - start) 1

private def rotateView (cells : List α) (first stop width count : Nat) : List α :=
  cells.take first ++ (sliceRotateRight
    ((cells.drop first).take (stop - first)) (stop - first) width count).2 ++ cells.drop stop

/-- Go 1.27.0 src/slices/slices.go:262–348. Within a common backing
    array, each element is represented by one address unit. A copy reads
    its entire source before overwriting the destination. -/
def sliceReplace (zero : α) (width : Nat) (s : GrowView α)
    (lo hi : Int) (source : InsertValues α) : Except Unit (GrowObservation α) :=
  let i := lo.toNat
  let j := hi.toNat
  let n := s.length
  if lo < 0 ∨ hi < lo ∨ n < j then .error ()
  else if i = j then sliceInsert s lo source
  else
    let values := source.read s.cells
    let m := values.length
    let suffix := (s.cells.drop j).take (n - j)
    let total := i + m + (n - j)
    if total > s.cells.length then
      .ok ⟨⟨s.cells.take i ++ values ++ suffix, total, false⟩, true⟩
    else
      let result := fun cells => .ok ⟨⟨cells, total, s.isNil⟩, false⟩
      let clearTail := fun cells => writeCells cells total (List.replicate (n - total) zero)
      if j = n then
        result (clearTail (writeCells s.cells i values))
      else if i + m ≤ j then
        let copied := writeCells s.cells i values
        let shifted := writeCells copied (i + m) ((copied.drop j).take (n - j))
        result (clearTail shifted)
      else if !overlapsTarget source (i + m) total then
        let shifted := writeCells s.cells (i + m) suffix
        result (writeCells shifted i (source.read shifted))
      else
        let y := m - (j - i)
        if !overlapsTarget source i j then
          let intoX := writeCells s.cells i (values.drop y)
          let intoY := writeCells intoX n ((source.read intoX).take y)
          result (rotateView intoY i total width y)
        else if !overlapsTarget source n total then
          let intoY := writeCells s.cells n (values.take y)
          let intoX := writeCells intoY i ((source.read intoY).drop y)
          result (rotateView intoX i total width y)
        else
          -- Overlap with both adjacent destination regions puts the entire
          -- untouched suffix inside the aliased source. Translate startIdx
          -- to a difference of aligned element offsets in one backing array.
          let k := match source with
            | .alias offset _ => j - offset
            | .independent _ => 0
          let copied := writeCells s.cells i values
          result (writeCells copied (i + m) ((copied.drop (i + k)).take (n - j)))

theorem sliceReplace_bounds (zero : α) (width : Nat) (s : GrowView α)
    (lo hi : Int) (source : InsertValues α)
    (bad : lo < 0 ∨ hi < lo ∨ s.length < hi.toNat) :
    sliceReplace zero width s lo hi source = .error () := by
  unfold sliceReplace
  rw [ite_eq_left bad]

theorem sliceReplace_insert (zero : α) (width : Nat) (s : GrowView α)
    (i : Nat) (source : InsertValues α) (valid : i ≤ s.length) :
    sliceReplace zero width s i i source = sliceInsert s i source := by
  have nonnegative : ¬ (i : Int) < 0 := by omega
  simp [sliceReplace, Nat.not_lt_of_ge valid, nonnegative]

theorem sliceReplace_allocate (zero : α) (width : Nat) (s : GrowView α)
    (i j : Nat) (source : InsertValues α)
    (nonempty : i < j) (endValid : j ≤ s.length)
    (insufficient : i + (source.read s.cells).length + (s.length - j) > s.cells.length) :
    sliceReplace zero width s i j source =
      .ok ⟨⟨s.cells.take i ++ source.read s.cells ++ (s.cells.take s.length).drop j,
        i + (source.read s.cells).length + (s.length - j), false⟩, true⟩ := by
  have suffix : (s.cells.drop j).take (s.length - j) = (s.cells.take s.length).drop j := by
    rw [List.take_drop, Nat.add_sub_cancel' endValid]
  simp [sliceReplace, Nat.not_lt_of_ge (Nat.le_of_lt nonempty), Nat.ne_of_lt nonempty,
    Nat.not_lt_of_ge endValid, insufficient, suffix]

theorem sliceReplace_clear (zero : α) (width : Nat) (visible spare values : List α)
    (nonempty : visible ≠ []) (shorter : values.length ≤ visible.length) :
    sliceReplace zero width ⟨visible ++ spare, visible.length, false⟩ 0 visible.length
      (.independent values) =
        .ok ⟨⟨values ++ List.replicate (visible.length - values.length) zero ++ spare,
          values.length, false⟩, false⟩ := by
  have unequal : 0 ≠ visible.length := by
    have positive : visible.length ≠ 0 := by simpa using nonempty
    omega
  have fits : ¬ visible.length + spare.length < values.length := by omega
  have nonnegative : ¬ (visible.length : Int) < 0 := by omega
  have afterValues : (values ++ (visible ++ spare).drop values.length).take values.length = values := by
    exact List.take_append_length
  have afterClear : (values ++ (visible ++ spare).drop values.length).drop
      (values.length + (visible.length - values.length)) = spare := by
    rw [← List.drop_drop, List.drop_append_length, List.drop_drop, Nat.add_sub_cancel' shorter]
    exact List.drop_append_length
  simp [sliceReplace, InsertValues.read, unequal, fits, nonnegative, writeCells, afterValues, afterClear,
    List.append_assoc]

end GoLean.Slices
