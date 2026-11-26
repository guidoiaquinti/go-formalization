namespace GoLean.Sort

/-- Go 1.27.0 src/sort/sort.go:122-125. Represent a slice by its visible
elements; both nil and non-nil empty slices have an empty sequence. -/
def intSliceLen (xs : List Int) : Nat := xs.length

theorem intSliceLen_subslice (xs : List Int) (lo hi : Nat)
    (ordered : lo ≤ hi) (bounded : hi ≤ xs.length) :
    intSliceLen ((xs.drop lo).take (hi - lo)) + lo = hi := by
  simp only [intSliceLen, List.length_take, List.length_drop]
  omega


end GoLean.Sort
