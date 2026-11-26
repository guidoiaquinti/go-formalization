namespace GoLean.Slices

private def compactTail (equal : α → α → Bool) (previous : α) : List α → List α
  | [] => []
  | current :: rest =>
    if equal current previous then compactTail equal current rest
    else current :: compactTail equal current rest

/-- Go 1.27.0 src/slices/slices.go:371–391. The two positions represent
the returned view and the original slice's backing region after zeroing. -/
def sliceCompact (equal : α → α → Bool) (zero : α) (xs : List α) : List α × List α :=
  let kept := match xs with
    | [] => []
    | first :: rest => first :: compactTail equal first rest
  (kept, kept ++ List.replicate (xs.length - kept.length) zero)

private theorem compactTail_spec (equal : α → α → Bool) (previous : α) (xs : List α) :
    compactTail equal previous xs = ((previous :: xs).zip xs).filterMap
      (fun (before, current) => if equal current before then none else some current) := by
  induction xs generalizing previous with
  | nil => simp [compactTail]
  | cons current rest ih =>
    simp only [compactTail, List.zip_cons_cons, List.filterMap_cons]
    rw [ih]
    cases equal current previous <;> rfl

private theorem compactTail_length (equal : α → α → Bool) (previous : α) (xs : List α) :
    (compactTail equal previous xs).length ≤ xs.length := by
  induction xs generalizing previous with
  | nil => simp [compactTail]
  | cons current rest ih =>
    have bound := ih current
    simp only [compactTail]
    split <;> simp only [List.length_cons] <;> omega

theorem sliceCompact_spec (equal : α → α → Bool) (zero : α) (first : α) (rest : List α) :
    let (kept, backing) := sliceCompact equal zero (first :: rest)
    kept = first :: ((first :: rest).zip rest).filterMap
      (fun (before, current) => if equal current before then none else some current) ∧
    backing.length = (first :: rest).length ∧
    backing.take kept.length = kept ∧
    backing.drop kept.length = List.replicate ((first :: rest).length - kept.length) zero := by
  have bound := compactTail_length equal first rest
  simp [sliceCompact, compactTail_spec equal first rest]
  rw [← compactTail_spec]
  omega

theorem sliceCompact_empty (equal : α → α → Bool) (zero : α) :
    sliceCompact equal zero [] = ([], []) := rfl

end GoLean.Slices
