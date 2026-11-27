import Init.Data.BitVec.Basic

namespace GoLean.Strings

structure AppendSliceState where
  bytes : Option (List (BitVec 8))
  capacity : Nat

structure AppendSliceWriteResult where
  bytes : Option (List (BitVec 8))
  capacity : Option Nat
  count : Nat
  error : Bool

/-- Go 1.27.0 src/strings/replace.go:301–306. The append reads the original
    source bytes, including when they alias the destination. None at the
    receiver denotes a nil pointer panic. A missing result capacity denotes
    growth, whose native capacity is at least the result length. -/
def appendSliceWriterWrite (receiver : Option AppendSliceState)
    (input : Option (List (BitVec 8))) : Option AppendSliceWriteResult := do
  let receiver ← receiver
  let source := input.getD []
  let bytes := receiver.bytes.getD [] ++ source
  let output := if receiver.bytes.isNone && source.isEmpty then none else some bytes
  let capacity := if bytes.length ≤ receiver.capacity then some receiver.capacity else none
  return ⟨output, capacity, source.length, false⟩

theorem appendSliceWriterWrite_result (receiver : AppendSliceState)
    (input : Option (List (BitVec 8))) :
    (appendSliceWriterWrite (some receiver) input).map (fun result => (result.count, result.error)) =
      some ((input.getD []).length, false) := rfl

theorem appendSliceWriterWrite_empty (receiver : AppendSliceState)
    (fits : (receiver.bytes.getD []).length ≤ receiver.capacity) :
    appendSliceWriterWrite (some receiver) (some []) =
      some ⟨receiver.bytes, some receiver.capacity, 0, false⟩ := by
  cases receiver with
  | mk bytes capacity =>
    cases bytes <;> simp_all [appendSliceWriterWrite]

theorem appendSliceWriterWrite_fits (bytes input : List (BitVec 8)) (capacity : Nat)
    (fits : bytes.length + input.length ≤ capacity) :
    appendSliceWriterWrite (some ⟨some bytes, capacity⟩) (some input) =
      some ⟨some (bytes ++ input), some capacity, input.length, false⟩ := by
  simp [appendSliceWriterWrite, fits]

theorem appendSliceWriterWrite_nilPointer (input : Option (List (BitVec 8))) :
    appendSliceWriterWrite none input = none := rfl

end GoLean.Strings
