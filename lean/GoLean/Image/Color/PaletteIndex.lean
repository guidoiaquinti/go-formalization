import GoLean.Image.Color.SqDiff

namespace GoLean.Image.Color

abbrev PaletteRGBA := BitVec 32 × BitVec 32 × BitVec 32 × BitVec 32
abbrev PaletteCall := Option (Except Nat PaletteRGBA)
abbrev PaletteResult := List Nat × Except Nat Nat

private def distance (a b : PaletteRGBA) : BitVec 32 :=
  sqDiff a.1 b.1 + sqDiff a.2.1 b.2.1 +
    sqDiff a.2.2.1 b.2.2.1 + sqDiff a.2.2.2 b.2.2.2

private def scan (query : PaletteRGBA) (i bestIndex : Nat) (best : BitVec 32) :
    List PaletteCall → PaletteResult
  | [] => ([], .ok bestIndex)
  | none :: _ => ([], .error 0)
  | some (.error payload) :: _ => ([i + 1], .error payload)
  | some (.ok channels) :: tail =>
    let score := distance query channels
    let (calls, result) :=
      if score < best then
        if score = 0 then ([], .ok i) else scan query (i + 1) i score tail
      else scan query (i + 1) bestIndex best tail
    ((i + 1) :: calls, result)

/-- The trace records the query as 0 and each reached callback as its index plus 1.
    A missing callback represents a nil Color interface. Go 1.27.0
    src/image/color/color.go:290-308. -/
def paletteIndex (query : PaletteCall) (palette : List PaletteCall) : PaletteResult :=
  match query with
  | none => ([], .error 0)
  | some (.error payload) => ([0], .error payload)
  | some (.ok channels) =>
    let (calls, result) := scan channels 0 0 0xffffffff palette
    (0 :: calls, result)

private def improve (state : Nat × Nat × BitVec 32) (score : BitVec 32) :
    Nat × Nat × BitVec 32 :=
  let (i, bestIndex, best) := state
  if score < best then (i + 1, i, score) else (i + 1, bestIndex, best)

private theorem improve_zero (i bestIndex : Nat) (scores : List (BitVec 32)) :
    (scores.foldl improve (i, bestIndex, 0)).2.1 = bestIndex := by
  induction scores generalizing i with
  | nil => rfl
  | cons score tail ih =>
    simp [List.foldl_cons, improve]
    exact ih (i + 1)

private theorem scan_return (query : PaletteRGBA) (i bestIndex : Nat) (best : BitVec 32)
    (palette : List PaletteRGBA) :
    (scan query i bestIndex best (palette.map fun c => some (.ok c))).2 =
      .ok (((palette.map (distance query)).foldl improve (i, bestIndex, best)).2.1) := by
  induction palette generalizing i bestIndex best with
  | nil => rfl
  | cons color tail ih =>
    simp only [List.map_cons, scan, List.foldl_cons, improve]
    split
    · split
      · rename_i zero
        simp [zero]
        exact (improve_zero (i + 1) i _).symm
      · cases h : scan query (i + 1) i (distance query color)
            (tail.map fun c => some (.ok c)) with
        | mk calls result => simpa [h] using ih (i + 1) i (distance query color)
    · cases h : scan query (i + 1) bestIndex best
          (tail.map fun c => some (.ok c)) with
      | mk calls result => simpa [h] using ih (i + 1) bestIndex best

/-- On normally returning callbacks, the return index equals a strict-minimum
    fold over all scores, even when the source stops at the first zero score. -/
theorem paletteIndex_return (query : PaletteRGBA) (palette : List PaletteRGBA) :
    (paletteIndex (some (.ok query)) (palette.map fun c => some (.ok c))).2 =
      .ok (((palette.map (distance query)).foldl improve (0, 0, 0xffffffff)).2.1) := by
  cases h : scan query 0 0 0xffffffff (palette.map fun c => some (.ok c)) with
  | mk calls result => simpa [paletteIndex, h] using scan_return query 0 0 0xffffffff palette

theorem paletteIndex_query (query : PaletteRGBA) (palette : List PaletteCall) (payload : Nat) :
    paletteIndex (some (.ok query)) [] = ([0], .ok 0) ∧
    paletteIndex none palette = ([], .error 0) ∧
    paletteIndex (some (.error payload)) palette = ([0], .error payload) :=
  ⟨rfl, rfl, rfl⟩

theorem paletteIndex_zero (query candidate : PaletteRGBA) (tail : List PaletteCall)
    (zero : distance query candidate = 0) :
    paletteIndex (some (.ok query)) (some (.ok candidate) :: tail) = ([0, 1], .ok 0) := by
  simp [paletteIndex, scan, zero]

end GoLean.Image.Color
