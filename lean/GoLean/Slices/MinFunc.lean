namespace GoLean.Slices

/-- Go 1.27.0 src/slices/sort.go:79–90, with a pure comparator. -/
def sliceMinFunc (cmp : α → α → Int) (xs : List α) : Except String α :=
  match xs with
  | [] => .error "slices.MinFunc: empty list"
  | x :: rest => .ok (rest.foldl (fun best next => if cmp next best < 0 then next else best) x)

private def chooseMinimum (key : α → Int) (best next : α) : α :=
  if key next < key best then next else best

private theorem foldMinimum_spec (key : α → Int) (rest : List α) (first : α) :
    let result := rest.foldl (chooseMinimum key) first
    (∀ x ∈ first :: rest, key result ≤ key x) ∧
      (first :: rest).find? (fun x => decide (key x ≤ key result)) = some result := by
  induction rest generalizing first with
  | nil => simp
  | cons next rest ih =>
    simp only [List.foldl_cons, chooseMinimum]
    by_cases order : key next < key first
    · simp only [ite_eq_left order]
      obtain ⟨bound, found⟩ := ih next
      have nextBound := bound next (by simp)
      have firstBound : key (rest.foldl (chooseMinimum key) next) < key first := by omega
      constructor
      · intro x member
        rcases List.mem_cons.mp member with rfl | member
        · exact Int.le_of_lt firstBound
        · exact bound x member
      · simpa [List.find?, Int.not_le.mpr firstBound] using found
    · simp only [ite_eq_right order]
      obtain ⟨bound, found⟩ := ih first
      have firstBound := bound first (by simp)
      constructor
      · intro x member
        simp only [List.mem_cons] at member
        rcases member with rfl | rfl | member
        · exact firstBound
        · omega
        · exact bound x (by simp [member])
      · by_cases firstMatches : key first ≤ key (rest.foldl (chooseMinimum key) first)
        · simpa [List.find?, firstMatches] using found
        · have nextMisses : ¬key next ≤ key (rest.foldl (chooseMinimum key) first) := by omega
          simpa [List.find?, firstMatches, nextMisses] using found

theorem sliceMinFunc_spec (cmp : α → α → Int) (key : α → Int) (xs : List α)
    (order : ∀ a b, cmp a b < 0 ↔ key a < key b) :
    match sliceMinFunc cmp xs with
    | .error message => xs = [] ∧ message = "slices.MinFunc: empty list"
    | .ok result =>
      (∀ x ∈ xs, key result ≤ key x) ∧
        xs.find? (fun x => decide (key x ≤ key result)) = some result := by
  cases xs with
  | nil => simp [sliceMinFunc]
  | cons first rest =>
    simp only [sliceMinFunc, order]
    exact foldMinimum_spec key rest first

end GoLean.Slices
