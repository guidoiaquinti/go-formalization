namespace GoLean.Sync

inductive OnceFault where
  | nilCallback
  | nilPanic
  | payload (code : Int)
  | nilReceiver
  deriving DecidableEq

inductive OnceCallback (α : Type) where
  | missing
  | raises (fault : OnceFault)
  | returns (value : α)

private def invoke (callback : OnceCallback α) : Bool × Except OnceFault α :=
  match callback with
  | .missing => (false, .error .nilCallback)
  | .raises fault => (true, .error fault)
  | .returns value => (true, .ok value)

structure OnceObservation where
  done : Option Bool
  calls : Bool
  result : Except OnceFault Unit

/-- Go 1.27.0 src/sync/once.go:73–80. The serial transition holds the
    mutex while calling f, and publishes done on both return and panic. -/
def onceDoSlow (done : Option Bool) (callback : OnceCallback Unit) : OnceObservation :=
  match done with
  | none => ⟨none, false, .error .nilReceiver⟩
  | some true => ⟨some true, false, .ok ()⟩
  | some false =>
    let (calls, result) := invoke callback
    ⟨some true, calls, result⟩

/-- Go 1.27.0 src/sync/once.go:52–71. The completed fast path bypasses f. -/
def onceDo (done : Option Bool) (callback : OnceCallback Unit) : OnceObservation :=
  if done = some true then ⟨some true, false, .ok ()⟩
  else onceDoSlow done callback

structure OnceClosure (α : Type) where
  pending : Option (OnceCallback α)
  cached : Option (Except OnceFault α)

/-- Go 1.27.0 src/sync/oncefunc.go:11–42. The created closure has no result yet. -/
def onceFunc (callback : OnceCallback Unit) : OnceClosure Unit := ⟨some callback, none⟩

/-- Go 1.27.0 src/sync/oncefunc.go:46–74. Its first successful value is cached. -/
def onceValue (callback : OnceCallback α) : OnceClosure α := ⟨some callback, none⟩

/-- Go 1.27.0 src/sync/oncefunc.go:80–109. Both values are cached together. -/
def onceValues (callback : OnceCallback (α × β)) : OnceClosure (α × β) :=
  ⟨some callback, none⟩

def callOnceClosure (s : OnceClosure α) : Bool × Except OnceFault α × OnceClosure α :=
  match s.pending with
  | some callback =>
    let (called, result) := invoke callback
    (called, result, ⟨none, some result⟩)
  | none =>
    (false, s.cached.getD (.error .nilCallback), s)

theorem onceDo_completed (callback : OnceCallback Unit) :
    onceDo (some true) callback = ⟨some true, false, .ok ()⟩ := rfl

theorem onceDoSlow_remembersPanic (payload : Int) :
    onceDoSlow (some false) (.raises (.payload payload)) =
      ⟨some true, true, .error (.payload payload)⟩ := rfl

theorem onceFunc_replaysPanic (payload : Int) :
    let (_, _, first) := callOnceClosure (onceFunc (.raises (.payload payload)))
    (callOnceClosure first).1 = false ∧
      (callOnceClosure first).2.1 = .error (.payload payload) := ⟨rfl, rfl⟩

theorem onceValue_caches (value : α) :
    let (_, _, first) := callOnceClosure (onceValue (.returns value))
    (callOnceClosure first).1 = false ∧
      (callOnceClosure first).2.1 = .ok value := ⟨rfl, rfl⟩

theorem onceValues_caches (x : α) (y : β) :
    let (_, _, first) := callOnceClosure (onceValues (.returns (x, y)))
    (callOnceClosure first).1 = false ∧
      (callOnceClosure first).2.1 = .ok (x, y) := ⟨rfl, rfl⟩

end GoLean.Sync
