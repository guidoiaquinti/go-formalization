import GoLean.Slices.Grow

namespace GoLean.Encoding.Binary

open GoLean.Slices

structure EnsureResult where
  buffer : List (BitVec 8)
  position : List (BitVec 8)
  nilBuffer : Bool
  allocated : Bool

def ensure (view : GrowView (BitVec 8)) (count : Int) : Except String EnsureResult :=
  (sliceGrow 0 view count).map fun grown =>
    let buffer := grown.view.cells.take (view.length + count.toNat)
    ⟨buffer, buffer.drop view.length, grown.view.isNil, grown.allocated⟩

theorem ensure_negative (view : GrowView (BitVec 8)) (count : Int) (negative : count < 0) :
    ensure view count = .error "cannot be negative" := by
  simp [ensure, Except.map, sliceGrow_negative _ _ _ negative]

theorem ensure_reuse (view : GrowView (BitVec 8)) (count : Nat)
    (fits : view.length + count ≤ view.cells.length) :
    ensure view count =
      .ok ⟨view.cells.take (view.length + count),
        (view.cells.take (view.length + count)).drop view.length, view.isNil, false⟩ := by
  simp [ensure, Except.map, sliceGrow_reuse _ _ _ fits]

theorem ensure_allocate (view : GrowView (BitVec 8)) (count : Nat)
    (valid : view.length ≤ view.cells.length) (short : view.cells.length < view.length + count) :
    ensure view count =
      let cells := view.cells ++ List.replicate (view.length + count - view.cells.length) 0
      let buffer := cells.take (view.length + count)
      .ok ⟨buffer, buffer.drop view.length, false, true⟩ := by
  simp [ensure, Except.map, sliceGrow_allocate _ _ _ valid short]

end GoLean.Encoding.Binary
