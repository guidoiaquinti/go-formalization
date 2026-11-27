namespace GoLean.Image.Color

/-- The identity is an abstract fresh allocation token, not a Go address. -/
structure ModelFuncWrapper (F : Type) where
  identity : Nat
  callback : Option F

/-- Go 1.27.0 src/image/color/color.go:147-157. Construct a pointer
    wrapper and store the callback, including nil, without evaluating it.
    The supply denotes the next unused identity. -/
def modelFunc (next : Nat) (callback : Option F) : ModelFuncWrapper F × Nat :=
  (⟨next, callback⟩, next + 1)

theorem modelFunc_spec (next : Nat) (callback : Option F) :
    (modelFunc next callback).1.callback = callback ∧
    (modelFunc next callback).1.identity = next ∧
    next < (modelFunc next callback).2 := by
  exact ⟨rfl, rfl, Nat.lt_succ_self next⟩

theorem modelFunc_distinct (next : Nat) (firstCallback secondCallback : Option F) :
    (modelFunc next firstCallback).1.identity ≠
      (modelFunc (modelFunc next firstCallback).2 secondCallback).1.identity := by
  exact Nat.ne_of_lt (Nat.lt_succ_self next)

end GoLean.Image.Color
