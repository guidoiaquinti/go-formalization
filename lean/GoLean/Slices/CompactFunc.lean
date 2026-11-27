import GoLean.Slices.Compact

namespace GoLean.Slices

private structure CompactScan (α : Type) where
  kept : List α
  trace : List (α × α)
  panic : Option Nat

private def compactFuncScan (equal : Option (Nat → α → α → Except Nat Bool))
    (i : Nat) (previous : α) : List α → CompactScan α
  | [] => ⟨[], [], none⟩
  | current :: rest =>
    match equal with
    | none => ⟨[], [], some 0⟩
    | some callback =>
      let pair := (current, previous)
      match callback i current previous with
      | .error payload => ⟨[], [pair], some payload⟩
      | .ok same =>
        let tail := compactFuncScan equal (i + 1) current rest
        ⟨if same then tail.kept else current :: tail.kept, pair :: tail.trace, tail.panic⟩

/-- Go 1.27.0 src/slices/slices.go:396–415. The scan visits original adjacent
    pairs. Each copy writes strictly before the current source cell and hence
    leaves the next comparison intact. A panic skips the final clear. -/
def sliceCompactFunc (equal : Option (Nat → α → α → Except Nat Bool)) (zero : α)
    (xs : List α) : Except Nat (List α) × List α × List (α × α) :=
  match xs with
  | [] => (.ok [], [], [])
  | first :: rest =>
    let scan := compactFuncScan equal 0 first rest
    let kept := first :: scan.kept
    match scan.panic with
    | none => (.ok kept, kept ++ List.replicate (xs.length - kept.length) zero, scan.trace)
    | some payload => (.error payload, kept ++ xs.drop kept.length, scan.trace)

private theorem compactFuncScan_pure (equal : α → α → Bool) (i : Nat)
    (previous : α) (rest : List α) :
    let pairs := rest.zip (previous :: rest)
    compactFuncScan (some (fun _ a b => .ok (equal a b))) i previous rest =
      ⟨pairs.filterMap (fun (a, b) => if equal a b then none else some a), pairs, none⟩ := by
  induction rest generalizing i previous with
  | nil => rfl
  | cons head tail ih =>
    cases h : equal head previous <;> simp [compactFuncScan, h, ih]

theorem sliceCompactFunc_pure (equal : α → α → Bool) (zero first : α) (rest : List α) :
    let pairs := rest.zip (first :: rest)
    let kept := first :: pairs.filterMap (fun (a, b) => if equal a b then none else some a)
    sliceCompactFunc (some (fun _ a b => .ok (equal a b))) zero (first :: rest) =
      (.ok kept, kept ++ List.replicate ((first :: rest).length - kept.length) zero, pairs) := by
  simp [sliceCompactFunc, compactFuncScan_pure]

private theorem pairOrder (equal : α → α → Bool) (first : α) (rest : List α) :
    ((rest.zip (first :: rest)).filterMap fun (a, b) => if equal a b then none else some a) =
    (((first :: rest).zip rest).filterMap fun (b, a) => if equal a b then none else some a) := by
  induction rest generalizing first with
  | nil => rfl
  | cons head tail ih =>
    simp only [List.zip_cons_cons, List.filterMap_cons, ih head]

theorem sliceCompactFunc_matchesCompact (equal : α → α → Bool) (zero first : α) (rest : List α) :
    let (kept, backing) := sliceCompact equal zero (first :: rest)
    sliceCompactFunc (some (fun _ a b => .ok (equal a b))) zero (first :: rest) =
      (.ok kept, backing, rest.zip (first :: rest)) := by
  have compact := sliceCompact_spec equal zero first rest
  cases h : sliceCompact equal zero (first :: rest) with
  | mk kept backing =>
    rw [h] at compact
    have writes := List.take_append_drop kept.length backing
    rw [compact.2.2.1, compact.2.2.2] at writes
    rw [sliceCompactFunc_pure, pairOrder, ← compact.1, writes]

theorem sliceCompactFunc_nil (zero first second : α) (rest : List α) :
    sliceCompactFunc none zero [] = (.ok [], [], []) ∧
    sliceCompactFunc none zero [first] = (.ok [first], [first], []) ∧
    sliceCompactFunc none zero (first :: second :: rest) =
      (.error 0, first :: second :: rest, []) :=
  ⟨rfl, rfl, rfl⟩

theorem sliceCompactFunc_partialWrite
    (equal : Nat → α → α → Except Nat Bool) (zero first second third fault : α)
    (tail : List α) (payload : Nat)
    (firstEqual : equal 0 second first = .ok true)
    (secondDistinct : equal 1 third second = .ok false)
    (thirdPanics : equal 2 fault third = .error payload) :
    sliceCompactFunc (some equal) zero (first :: second :: third :: fault :: tail) =
      (.error payload, first :: third :: third :: fault :: tail,
        [(second, first), (third, second), (fault, third)]) := by
  simp [sliceCompactFunc, compactFuncScan, firstEqual, secondDistinct, thirdPanics]

end GoLean.Slices
