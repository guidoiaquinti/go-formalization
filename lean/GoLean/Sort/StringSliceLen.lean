namespace GoLean.Sort

/-- Go 1.27.0 src/sort/sort.go:155-158. Represent a slice by its visible
elements; both nil and non-nil empty slices have an empty sequence. -/
def stringSliceLen (xs : List ByteArray) : Nat := xs.length

theorem stringSliceLen_subslice (xs : List ByteArray) (lo hi : Nat)
    (ordered : lo ≤ hi) (bounded : hi ≤ xs.length) :
    stringSliceLen ((xs.drop lo).take (hi - lo)) + lo = hi := by
  simp only [stringSliceLen, List.length_take, List.length_drop]
  omega


end GoLean.Sort
