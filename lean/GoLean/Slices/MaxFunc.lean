import GoLean.Slices.MinFunc

namespace GoLean.Slices

/-- Go 1.27.0 src/slices/sort.go:109–120. Reversing the sign of the
mathematical comparator result expresses Go's positive branch as a negative test. -/
def sliceMaxFunc (cmp : α → α → Int) (xs : List α) : Except String α :=
  (sliceMinFunc (fun a b => -cmp a b) xs).mapError (fun _ => "slices.MaxFunc: empty list")

theorem sliceMaxFunc_spec (cmp : α → α → Int) (key : α → Int) (xs : List α)
    (order : ∀ a b, 0 < cmp a b ↔ key b < key a) :
    match sliceMaxFunc cmp xs with
    | .error message => xs = [] ∧ message = "slices.MaxFunc: empty list"
    | .ok result =>
      (∀ x ∈ xs, key x ≤ key result) ∧
        xs.find? (fun x => decide (key result ≤ key x)) = some result := by
  have reverseOrder (a b : α) : -cmp a b < 0 ↔ -key a < -key b := by
    simpa using order a b
  cases xs with
  | nil => simp [sliceMaxFunc, sliceMinFunc, Except.mapError]
  | cons first rest =>
    have spec := sliceMinFunc_spec (fun a b => -cmp a b) (fun a => -key a) (first :: rest) reverseOrder
    simpa [sliceMaxFunc, sliceMinFunc, Except.mapError] using spec

end GoLean.Slices
