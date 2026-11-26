namespace GoLean.Sort

/-- Go 1.27.0 src/sort/sort.go:132-136. Represent a slice by its visible
binary64 bit patterns; both nil and non-nil empty slices have an empty sequence. -/
def float64SliceLen (xs : List UInt64) : Nat := xs.length

theorem float64SliceLen_subslice (xs : List UInt64) (lo hi : Nat)
    (ordered : lo ≤ hi) (bounded : hi ≤ xs.length) :
    float64SliceLen ((xs.drop lo).take (hi - lo)) + lo = hi := by
  simp only [float64SliceLen, List.length_take, List.length_drop]
  omega


end GoLean.Sort
