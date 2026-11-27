import GoLean.Strings.AppendSliceWriterWrite

namespace GoLean.Strings

/-- Go 1.27.0 src/strings/replace.go:308–312. Append the string's raw
    bytes and return its byte length, with the same nil and capacity
    behavior as the byte writer. -/
def appendSliceWriterWriteString (receiver : Option AppendSliceState)
    (input : List (BitVec 8)) : Option AppendSliceWriteResult :=
  appendSliceWriterWrite receiver (some input)

theorem appendSliceWriterWriteString_result (receiver : AppendSliceState)
    (input : List (BitVec 8)) :
    (appendSliceWriterWriteString (some receiver) input).map (fun result => (result.count, result.error)) =
      some (input.length, false) :=
  appendSliceWriterWrite_result receiver (some input)

theorem appendSliceWriterWriteString_empty (receiver : AppendSliceState)
    (fits : (receiver.bytes.getD []).length ≤ receiver.capacity) :
    appendSliceWriterWriteString (some receiver) [] =
      some ⟨receiver.bytes, some receiver.capacity, 0, false⟩ :=
  appendSliceWriterWrite_empty receiver fits

theorem appendSliceWriterWriteString_fits (bytes input : List (BitVec 8)) (capacity : Nat)
    (fits : bytes.length + input.length ≤ capacity) :
    appendSliceWriterWriteString (some ⟨some bytes, capacity⟩) input =
      some ⟨some (bytes ++ input), some capacity, input.length, false⟩ :=
  appendSliceWriterWrite_fits bytes input capacity fits

theorem appendSliceWriterWriteString_nilPointer (input : List (BitVec 8)) :
    appendSliceWriterWriteString none input = none := rfl

end GoLean.Strings
