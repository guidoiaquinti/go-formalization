namespace GoLean.Sort

structure SortInterface (Length Test Swap : Type) where
  length : Unit → Length
  less : Int → Int → Test
  swap : Int → Int → Swap

/-- Go 1.27.0 src/sort/sort.go:90-104 embeds the interface and overrides Less.
Callback outcomes can include state transformations and panics. -/
def reverseInterface (data : SortInterface Length Test Swap) : SortInterface Length Test Swap :=
  { data with less := fun i j => data.less j i }

theorem reverseInterface_spec (data : SortInterface Length Test Swap) (i j : Int) :
    (reverseInterface data).length () = data.length () ∧
    (reverseInterface data).less i j = data.less j i ∧
    (reverseInterface data).swap i j = data.swap i j :=
  ⟨rfl, rfl, rfl⟩


end GoLean.Sort
