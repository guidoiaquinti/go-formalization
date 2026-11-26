namespace GoLean.Slices

/-- Go 1.27.0 src/slices/slices.go:482-487. The list model exchanges
the endpoints and continues on the interior, which is the range between
the updated loop indices. -/
def sliceReverse (xs : List α) : List α :=
  match xs with
  | [] => []
  | first :: rest =>
    if h : rest = [] then [first]
    else rest.getLast h :: (sliceReverse rest.dropLast ++ [first])
termination_by xs.length
decreasing_by simp; omega

theorem sliceReverse_spec (xs : List α) : sliceReverse xs = xs.reverse := by
  cases xs with
  | nil => rw [sliceReverse]; rfl
  | cons first rest =>
    rw [sliceReverse]
    split
    · rename_i empty
      simp [empty]
    · rename_i nonempty
      rw [sliceReverse_spec rest.dropLast, List.reverse_cons]
      have tailReverse :
          rest.reverse = rest.getLast nonempty :: rest.dropLast.reverse := by
        simpa [List.reverse_append] using
          (congrArg List.reverse (List.dropLast_concat_getLast nonempty)).symm
      rw [tailReverse]
      rfl
termination_by xs.length
decreasing_by simp; omega

end GoLean.Slices
