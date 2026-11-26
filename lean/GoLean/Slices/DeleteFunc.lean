namespace GoLean.Slices

private structure DeleteScan (α : Type) where
  panicked : Bool
  kept : List α
  calls : Nat

private def deleteScan (predicate : α → Except Unit Bool) : List α → DeleteScan α
  | [] => ⟨false, [], 0⟩
  | value :: rest =>
    match predicate value with
    | .error () => ⟨true, [], 1⟩
    | .ok remove =>
      let scan := deleteScan predicate rest
      ⟨scan.panicked, if remove then scan.kept else value :: scan.kept, scan.calls + 1⟩

/-- Go 1.27.0 src/slices/slices.go:241–255 and 105–114. The written
prefix comprises the retained items visited before termination. On a callback
panic the unwritten suffix stays intact and the final clear is skipped. -/
def sliceDeleteFunc (predicate : α → Except Unit Bool) (zero : α)
    (xs : List α) : Option (List α) × List α × List α :=
  let scan := deleteScan predicate xs
  let trace := xs.take scan.calls
  if scan.panicked then (none, scan.kept ++ xs.drop scan.kept.length, trace)
  else (some scan.kept, scan.kept ++ List.replicate (xs.length - scan.kept.length) zero, trace)

private theorem deleteScan_pure (remove : α → Bool) (xs : List α) :
    deleteScan (fun x => .ok (remove x)) xs =
      ⟨false, xs.filter fun x => !remove x, xs.length⟩ := by
  induction xs with
  | nil => rfl
  | cons x xs ih =>
    cases h : remove x <;> simp [deleteScan, ih, h]

theorem sliceDeleteFunc_pure (remove : α → Bool) (zero : α) (xs : List α) :
    let kept := xs.filter fun x => !remove x
    sliceDeleteFunc (fun x => .ok (remove x)) zero xs =
      (some kept, kept ++ List.replicate (xs.length - kept.length) zero, xs) := by
  simp [sliceDeleteFunc, deleteScan_pure]

private theorem deleteScan_panic (predicate : α → Except Unit Bool) (remove : α → Bool)
    (before : List α) (fault : α) (after : List α)
    (prefixOK : ∀ x ∈ before, predicate x = .ok (remove x))
    (faultError : predicate fault = .error ()) :
    deleteScan predicate (before ++ fault :: after) =
      ⟨true, before.filter fun x => !remove x, before.length + 1⟩ := by
  induction before with
  | nil => simp [deleteScan, faultError]
  | cons x xs ih =>
    have headOK := prefixOK x (by simp)
    have tailOK : ∀ y ∈ xs, predicate y = .ok (remove y) :=
      fun y member => prefixOK y (List.mem_cons_of_mem x member)
    cases h : remove x <;> simp [deleteScan, headOK, ih tailOK, h]

theorem sliceDeleteFunc_panic (predicate : α → Except Unit Bool) (remove : α → Bool)
    (zero : α) (before : List α) (fault : α) (after : List α)
    (prefixOK : ∀ x ∈ before, predicate x = .ok (remove x))
    (faultError : predicate fault = .error ()) :
    let kept := before.filter fun x => !remove x
    let xs := before ++ fault :: after
    sliceDeleteFunc predicate zero xs =
      (none, kept ++ xs.drop kept.length, before ++ [fault]) := by
  simp [sliceDeleteFunc, deleteScan_panic predicate remove before fault after prefixOK faultError,
    List.take_append, List.take_of_length_le (Nat.le_succ _)]

end GoLean.Slices
