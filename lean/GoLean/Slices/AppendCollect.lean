namespace GoLean.Slices

inductive AppendSeqError (ε : Type) where
  | nilSequence
  | producer (payload : ε)

structure AppendSeqObservation (α ε : Type) where
  result : Except (AppendSeqError ε) (Option (List α))
  oldBacking : List α
  acknowledged : Nat
  reused : Bool

/-- Go 1.27.0 src/slices/iter.go:50–55. The producer supplies the values
    yielded before returning or panicking. Appends overwrite the old spare
    slots until they are exhausted, then use newly allocated backing.
    Allocator-rounded capacity is outside this observation. -/
def sliceAppendSeq (initial : Option (List α)) (spare : List α)
    (sequence : Option (List α × Except ε Unit)) : AppendSeqObservation α ε :=
  let initialValues := initial.getD []
  match sequence with
  | none => ⟨.error .nilSequence, initialValues ++ spare, 0, false⟩
  | some (values, producer) =>
    let writes := values.take spare.length
    let old := initialValues ++ writes ++ spare.drop writes.length
    let output := if values.isEmpty then initial else some (initialValues ++ values)
    let reused := match producer with
      | .ok _ => !(initialValues ++ spare).isEmpty && decide (values.length ≤ spare.length)
      | .error _ => false
    ⟨(producer.map fun _ => output).mapError AppendSeqError.producer, old, values.length, reused⟩

/-- Go 1.27.0 src/slices/iter.go:59–61. Collect starts from a nil slice. -/
def sliceCollect (sequence : Option (List α × Except ε Unit)) : AppendSeqObservation α ε :=
  sliceAppendSeq none [] sequence

theorem sliceAppendSeq_empty (initial : Option (List α)) (spare : List α) :
    (sliceAppendSeq (ε := ε) initial spare (some ([], .ok ()))).result = .ok initial := rfl

theorem sliceAppendSeq_nilSequence (initial : Option (List α)) (spare : List α) :
    sliceAppendSeq (ε := ε) initial spare none =
      ⟨.error .nilSequence, initial.getD [] ++ spare, 0, false⟩ := rfl

theorem sliceAppendSeq_producerPanic (payload : ε) :
    sliceAppendSeq (some [11]) [90, 91] (some ([4, 5, 6], .error payload)) =
      ⟨.error (.producer payload), [11, 4, 5], 3, false⟩ := rfl

theorem sliceCollect_empty :
    (sliceCollect (α := α) (ε := ε) (some ([], .ok ()))).result = .ok none := rfl

theorem sliceCollect_values (first : α) (remaining : List α) :
    sliceCollect (ε := ε) (some (first :: remaining, .ok ())) =
      ⟨.ok (some (first :: remaining)), [], remaining.length + 1, false⟩ := rfl

theorem sliceCollect_producerPanic (values : List α) (payload : ε) :
    (sliceCollect (some (values, .error payload))).result = .error (.producer payload) := rfl

end GoLean.Slices
