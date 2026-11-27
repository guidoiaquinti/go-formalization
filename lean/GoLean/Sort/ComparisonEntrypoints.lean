import GoLean.Sort.IsSorted

namespace GoLean.Sort

inductive ReflectedLength where
  | supported
  | invalid
  | unsupported

/-- Go 1.27.0 src/sort/slice.go:52. Len validates the reflection shape
before the reverse adjacent scan invokes the supplied comparison. -/
def sliceIsSorted (shape : ReflectedLength) (n : Nat)
    (less : Nat → Nat → Except Unit Bool) : Except String Bool × List Nat :=
  match shape with
  | .invalid => (.error "invalid", [])
  | .unsupported => (.error "unsupported", [])
  | .supported =>
    let (result, calls) := isSorted n less
    (result.mapError (fun _ => "callback"), calls)

/-- Go 1.27.0 src/sort/sort.go:97. The receiver dispatches a single call
on the embedded interface with the original signed indices reversed. -/
def reverseLess (less : Option (Int → Int → Except Unit Bool)) (i j : Int) :
    Except String Bool × List (Int × Int) :=
  match less with
  | none => (.error "nil", [])
  | some call => ((call j i).mapError (fun _ => "callback"), [(j, i)])

theorem sliceIsSorted_valid (n : Nat) (less : Nat → Nat → Except Unit Bool) :
    (sliceIsSorted .supported n less).1 = .ok true ↔
      ∀ i < n - 1, less (i + 1) i = .ok false := by
  have normalized (result : Except Unit Bool) :
      result.mapError (fun _ => "callback") = .ok true ↔ result = .ok true := by
    cases result <;> simp [Except.mapError]
  simpa [sliceIsSorted, normalized] using isSorted_spec n less

theorem sliceIsSorted_invalid (n : Nat) (less : Nat → Nat → Except Unit Bool) :
    sliceIsSorted .invalid n less = (.error "invalid", []) := rfl

theorem reverseLess_order (less : Int → Int → Except Unit Bool) (i j : Int) :
    reverseLess (some less) i j = ((less j i).mapError (fun _ => "callback"), [(j, i)]) := rfl

theorem reverseLess_nil (i j : Int) :
    reverseLess none i j = (.error "nil", []) := rfl

end GoLean.Sort
