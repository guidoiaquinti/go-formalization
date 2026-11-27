import GoLean.Strings.BuilderWriteString

namespace GoLean.Strings

/-- Go 1.27.0 strings/builder.go:87–91. Write uses the same copy check,
    raw byte append and result as WriteString. None represents a nil slice. -/
def builderWrite (receiver : Option BuilderState) (copied : Bool)
    (input : Option (List (BitVec 8))) : Option BuilderWriteResult :=
  builderWriteString receiver copied (input.getD [])

theorem builderWrite_success (b : BuilderState) (copied : Bool)
    (input : Option (List (BitVec 8))) (canWrite : (b.identitySet && copied) = false) :
    builderWrite (some b) copied input =
      some ⟨(some (input.getD []).length), (b.buffer.getD ([], 0)).1 ++ input.getD [], true⟩ :=
  builderWriteString_success b copied (input.getD []) canWrite

theorem builderWrite_copied (b : BuilderState) (input : Option (List (BitVec 8)))
    (identified : b.identitySet = true) :
    builderWrite (some b) true input =
      some ⟨none, (b.buffer.getD ([], 0)).1, true⟩ :=
  builderWriteString_copied b (input.getD []) identified

theorem builderWrite_nil (copied : Bool) (input : Option (List (BitVec 8))) :
    builderWrite none copied input = none := rfl

theorem builderWrite_nilSlice (receiver : Option BuilderState) (copied : Bool) :
    builderWrite receiver copied none = builderWrite receiver copied (some []) := rfl

end GoLean.Strings
