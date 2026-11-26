namespace GoLean.Slices

/-- Go 1.27.0 src/slices/slices.go:489-506 first checks the cumulative length,
then appends each slice in order. The length bound is Go's MaxInt. -/
def sliceConcat (maxLength : Nat) (parts : List (List α)) : Except String (Option (List α)) :=
  let size := (parts.map List.length).sum
  if size > maxLength then .error "len out of range"
  else if size = 0 then .ok none
  else .ok (some (parts.foldl (· ++ ·) []))

private theorem appendParts (parts : List (List α)) (initial : List α) :
    parts.foldl (· ++ ·) initial = initial ++ parts.flatten := by
  induction parts generalizing initial with
  | nil => simp
  | cons part rest ih => simp [ih, List.append_assoc]

theorem sliceConcat_spec (maxLength : Nat) (parts : List (List α)) :
    sliceConcat maxLength parts =
      if parts.flatten.length > maxLength then .error "len out of range"
      else if parts.flatten = [] then .ok none
      else .ok (some parts.flatten) := by
  simp only [sliceConcat, ← List.length_flatten, List.length_eq_zero_iff, appendParts,
    List.nil_append]

end GoLean.Slices
