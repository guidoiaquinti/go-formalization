namespace GoLean.Slices

/-- Go 1.27.0 src/slices/sort.go:92-104. The binary operation represents
the typed Go max builtin. -/
def sliceMax (maximum : α → α → α) (xs : List α) : Except String α :=
  match xs with
  | [] => .error "slices.Max: empty list"
  | x :: rest => .ok (rest.foldl maximum x)

private theorem foldMax_spec (xs : List Int) (initial : Int) :
    let result := xs.foldl max initial
    result ∈ initial :: xs ∧ ∀ x ∈ initial :: xs, x ≤ result := by
  induction xs generalizing initial with
  | nil => simp
  | cons x xs ih =>
    simp only [List.foldl_cons]
    have h := ih (max initial x)
    rw [Int.max_def] at h ⊢
    by_cases order : initial ≤ x <;> simp_all [List.mem_cons] <;> grind

theorem sliceMax_int_spec (xs : List Int) :
    match sliceMax max xs with
    | .error message => xs = [] ∧ message = "slices.Max: empty list"
    | .ok value => value ∈ xs ∧ ∀ x ∈ xs, x ≤ value := by
  cases xs with
  | nil => simp [sliceMax]
  | cons x xs => exact foldMax_spec xs x

end GoLean.Slices
