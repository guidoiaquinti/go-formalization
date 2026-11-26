namespace GoLean.Slices

/-- Go 1.27.0 src/slices/slices.go:351-363. None represents nil.
A non-nil slice copies its elements into a fresh slice. The model represents
the shallow copy by its element sequence. -/
def sliceClone (s : Option (List α)) : Option (List α) :=
  match s with
  | none => none
  | some xs => some ([] ++ xs)

theorem sliceClone_spec (s : Option (List α)) : sliceClone s = s := by
  cases s <;> rfl


end GoLean.Slices
