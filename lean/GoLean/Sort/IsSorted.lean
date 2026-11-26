namespace GoLean.Sort

/-- Go 1.27.0 src/sort/sort.go:110-118. The bound is the number of
adjacent pairs after Len has returned a nonnegative slice length. -/
private def checkDescending (less : Nat → Nat → Except Unit Bool) :
    Nat → Except Unit Bool × List Nat
  | 0 => (.ok true, [])
  | i + 1 =>
    match less (i + 1) i with
    | .error () => (.error (), [i + 1])
    | .ok true => (.ok false, [i + 1])
    | .ok false =>
      let (result, calls) := checkDescending less i
      (result, (i + 1) :: calls)

def isSorted (length : Nat) (less : Nat → Nat → Except Unit Bool) :
    Except Unit Bool × List Nat :=
  checkDescending less (length - 1)

private theorem checkDescending_spec (less : Nat → Nat → Except Unit Bool) (n : Nat) :
    (checkDescending less n).1 = .ok true ↔
      ∀ i < n, less (i + 1) i = .ok false := by
  induction n with
  | zero => simp [checkDescending]
  | succ n ih =>
    have step :
        (∀ i < n + 1, less (i + 1) i = .ok false) ↔
          less (n + 1) n = .ok false ∧ ∀ i < n, less (i + 1) i = .ok false := by
      constructor
      · intro all
        exact ⟨all n (by omega), fun i h => all i (by omega)⟩
      · rintro ⟨last, earlier⟩ i hi
        by_cases h : i = n
        · simpa [h] using last
        · exact earlier i (by omega)
    rw [step]
    cases h : less (n + 1) n with
    | error e => cases e; simp [checkDescending, h]
    | ok b => cases b <;> simp [checkDescending, h, ih]

theorem isSorted_spec (length : Nat) (less : Nat → Nat → Except Unit Bool) :
    (isSorted length less).1 = .ok true ↔
      ∀ i < length - 1, less (i + 1) i = .ok false :=
  checkDescending_spec less (length - 1)

/-- Apply the descending adjacent scan to a pure element comparison. -/
def isSortedList (less : α → α → Bool) (fallback : α) (xs : List α) : Bool :=
  match (isSorted xs.length (fun i j => .ok (less (xs[i]?.getD fallback) (xs[j]?.getD fallback)))).1 with
  | .ok result => result
  | .error () => false

theorem isSortedList_spec (less : α → α → Bool) (fallback : α) (xs : List α) :
    isSortedList less fallback xs = true ↔
      ∀ i < xs.length - 1, less (xs[i + 1]?.getD fallback) (xs[i]?.getD fallback) = false := by
  have resultTrue (result : Except Unit Bool) :
      (match result with | .ok b => b | .error () => false) = true ↔ result = .ok true := by
    cases result with
    | error e => cases e; simp
    | ok b => simp
  simp [isSortedList, resultTrue, isSorted_spec]

end GoLean.Sort
