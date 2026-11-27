namespace GoLean.Image.Color

inductive ConvertPanic (ε : Type) where
  | nilReceiver
  | nilCallback
  | callback (value : ε)

/-- The callback can change state before returning or panicking. The optional
    receiver and callback model the two distinct missing values.
    https://github.com/golang/go/blob/go1.27.0/src/image/color/color.go#L156-L162 -/
def modelFuncConvert (m : Option (Option (α → σ → σ × Except ε α)))
    (c : α) (s : σ) : List α × σ × Except (ConvertPanic ε) α :=
  match m with
  | none => ([], s, .error .nilReceiver)
  | some none => ([], s, .error .nilCallback)
  | some (some callback) =>
    let (updated, result) := callback c s
    ([c], updated, result.mapError ConvertPanic.callback)

theorem modelFuncConvert_spec (f : α → σ → σ × Except ε α)
    (c : α) (s : σ) :
    modelFuncConvert (some (some f)) c s =
      ([c], (f c s).1, match (f c s).2 with
        | .ok value => .ok value
        | .error value => .error (.callback value)) := by
  unfold modelFuncConvert
  cases h : f c s with
  | mk state result => cases result <;> simp [h, Except.mapError]

theorem modelFuncConvert_nilReceiver (c : α) (s : σ) :
    modelFuncConvert (ε := ε) none c s = ([], s, .error .nilReceiver) := rfl

theorem modelFuncConvert_nilCallback (c : α) (s : σ) :
    modelFuncConvert (ε := ε) (some none) c s = ([], s, .error .nilCallback) := rfl

end GoLean.Image.Color
