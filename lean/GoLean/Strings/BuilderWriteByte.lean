import GoLean.Strings.BuilderWriteString

namespace GoLean.Strings

/-- WriteByte uses the same copy check and append as WriteString,
    specialized to one byte. Its success has no returned count.
    https://github.com/golang/go/blob/go1.27.0/src/strings/builder.go#L32-L42
    https://github.com/golang/go/blob/go1.27.0/src/strings/builder.go#L94-L99 -/
def builderWriteByte (receiver : Option BuilderState) (copied : Bool)
    (byte : BitVec 8) : Option (Bool × List (BitVec 8) × Bool) :=
  (builderWriteString receiver copied [byte]).map fun r =>
    (r.count.isSome, r.contents, r.identitySet)

theorem builderWriteByte_success (b : BuilderState) (copied : Bool)
    (byte : BitVec 8) (canWrite : (b.identitySet && copied) = false) :
    builderWriteByte (some b) copied byte =
      some (true, (b.buffer.getD ([], 0)).1 ++ [byte], true) := by
  rw [builderWriteByte, builderWriteString_success b copied [byte] canWrite]
  rfl

theorem builderWriteByte_copied (b : BuilderState) (byte : BitVec 8)
    (identified : b.identitySet = true) :
    builderWriteByte (some b) true byte =
      some (false, (b.buffer.getD ([], 0)).1, true) := by
  rw [builderWriteByte, builderWriteString_copied b [byte] identified]
  rfl

theorem builderWriteByte_nil (copied : Bool) (byte : BitVec 8) :
    builderWriteByte none copied byte = none := rfl

end GoLean.Strings
