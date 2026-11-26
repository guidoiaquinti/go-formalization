namespace GoLean.Slices

/-- Go 1.27.0 src/slices/sort.go:62-74. The binary operation represents
the typed Go min builtin. -/
def sliceMin (minimum : α → α → α) (xs : List α) : Except String α :=
  match xs with
  | [] => .error "slices.Min: empty list"
  | x :: rest => .ok (rest.foldl minimum x)

private theorem foldMin_spec (xs : List Int) (initial : Int) :
    let result := xs.foldl min initial
    result ∈ initial :: xs ∧ ∀ x ∈ initial :: xs, result ≤ x := by
  induction xs generalizing initial with
  | nil => simp
  | cons x xs ih =>
    simp only [List.foldl_cons]
    have h := ih (min initial x)
    rw [Int.min_def] at h ⊢
    by_cases order : initial ≤ x <;> simp_all [List.mem_cons] <;> grind

theorem sliceMin_int_spec (xs : List Int) :
    match sliceMin min xs with
    | .error message => xs = [] ∧ message = "slices.Min: empty list"
    | .ok value => value ∈ xs ∧ ∀ x ∈ xs, value ≤ x := by
  cases xs with
  | nil => simp [sliceMin]
  | cons x xs => exact foldMin_spec xs x

end GoLean.Slices
