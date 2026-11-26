import GoLean.Sort.IndexedLess
import Init.Data.Float.Float

namespace GoLean.Sort

/-- Go 1.27.0 src/sort/sort.go:138-150. The second comparison uses
self-inequality to identify NaNs after the ordinary comparison. -/
def float64SliceLess (xs : List Float) (i j : Int) : Except Int Bool :=
  indexedLess (fun a b => a.lt b || (a != a && !(b != b))) xs i j

theorem float64SliceLess_spec (xs : List Float) (i j : Int) :
    float64SliceLess xs i j =
      if i < 0 ∨ xs.length ≤ i.toNat then .error i
      else if j < 0 ∨ xs.length ≤ j.toNat then .error j
      else
        let a := xs[i.toNat]?.getD 0
        let b := xs[j.toNat]?.getD 0
        .ok (a.lt b || (!(a == a) && (b == b))) := by
  have h := indexedLess_spec (fun (a b : Float) => a.lt b || (a != a && !(b != b))) 0 xs i j
  simpa [float64SliceLess, bne] using h

end GoLean.Sort
