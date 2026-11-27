import GoLean.Image.Color.PaletteIndex

namespace GoLean.Image.Color

/-- Each entry pairs the callback observed by Index with the interface value
    present in that slot when Convert reads it. The latter may have been
    replaced during a callback. None denotes a nil Color interface.
    Go 1.27.0 src/image/color/color.go:279-288. -/
def paletteConvert (query : PaletteCall) (palette : List (PaletteCall × Option α)) :
    List Nat × Except Nat (Option α) :=
  match palette with
  | [] => ([], .ok none)
  | _ :: _ =>
    let (calls, index) := paletteIndex query (palette.map Prod.fst)
    (calls, index.map fun i => palette[i]?.bind Prod.snd)

theorem paletteConvert_empty (query : PaletteCall) :
    paletteConvert (α := α) query [] = ([], .ok none) := rfl

theorem paletteConvert_selected (query : PaletteCall)
    (palette : List (PaletteCall × Option α)) (index : Nat) (calls : List Nat)
    (entry : PaletteCall × Option α) (atIndex : palette[index]? = some entry)
    (selection : paletteIndex query (palette.map Prod.fst) = (calls, .ok index)) :
    paletteConvert query palette = (calls, .ok entry.2) := by
  cases palette with
  | nil => simp at atIndex
  | cons head tail =>
    unfold paletteConvert
    rw [selection]
    simp [Except.map, atIndex]

theorem paletteConvert_panic (query : PaletteCall) (head : PaletteCall × Option α)
    (tail : List (PaletteCall × Option α)) (calls : List Nat) (payload : Nat)
    (selection : paletteIndex query ((head :: tail).map Prod.fst) = (calls, .error payload)) :
    paletteConvert query (head :: tail) = (calls, .error payload) := by
  unfold paletteConvert
  rw [selection]
  rfl

end GoLean.Image.Color
