import GoLean.Strings.StringWriterWriteString

namespace GoLean.Strings

structure StringWriterSelection where
  original : Bool
  writeString : List (BitVec 8) → Int × Nat

/-- Go 1.27.0 src/strings/replace.go:322–329. The existing StringWriter
    method survives the interface assertion unchanged. Otherwise expose
    the ordinary byte writer through the adapter. The option denotes
    method capability, independent of whether the dynamic receiver is nil. -/
def getStringWriter (write : List (BitVec 8) → Int × Nat)
    (writeString : Option (List (BitVec 8) → Int × Nat)) : StringWriterSelection :=
  match writeString with
  | some method => ⟨true, method⟩
  | none => ⟨false, fun input =>
      let result := stringWriterWriteString write input
      (result.count, result.error)⟩

theorem getStringWriter_original (write method : List (BitVec 8) → Int × Nat) :
    (getStringWriter write (some method)).original = true ∧
    (getStringWriter write (some method)).writeString = method :=
  ⟨rfl, rfl⟩

theorem getStringWriter_adapter (write : List (BitVec 8) → Int × Nat) (input : List (BitVec 8)) :
    (getStringWriter write none).original = false ∧
    (getStringWriter write none).writeString input = write input := by
  exact ⟨rfl, Prod.ext (stringWriterWriteString_count write input)
    (stringWriterWriteString_error write input)⟩

end GoLean.Strings
