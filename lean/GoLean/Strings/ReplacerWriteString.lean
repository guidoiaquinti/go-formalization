import GoLean.Strings.ReplacerReplace
import GoLean.Strings.ByteReplacerWriteString
import GoLean.Strings.ByteStringReplacerWriteString

namespace GoLean.Strings

private def writeBuilt (implementation : BuiltReplacer) (input : List (BitVec 8))
    (stringWriter : Bool) (script : List ReplacerWriteReply) : Option ReplacerWriteResult :=
  match implementation with
  | .single state =>
    some (singleStringReplacerWriteString state.finder.pattern state.value input stringWriter script)
  | .generic pairs _ => some (genericReplacerWriteString pairs input stringWriter script)
  | .bytes lookup => byteReplacerWriteString (some lookup) input stringWriter script
  | .byteStrings state => byteStringReplacerWriteString (some state.lookup) input stringWriter script

/-- Go 1.27.0 src/strings/replace.go:100–104. Lazily initialize, then
    use the selected implementation's own writer dispatch and error handling. -/
def replacerWriteString (receiver : Option ReplacerState) (input : List (BitVec 8))
    (stringWriter : Bool) (script : List ReplacerWriteReply) :
    Option (ReplacerState × ReplacerWriteResult) := do
  let state ← receiver
  let implementation := replacerImplementation state
  let result ← writeBuilt implementation input stringWriter script
  return (.ready implementation, result)

theorem replacerWriteString_nil (input : List (BitVec 8)) (mode : Bool)
    (script : List ReplacerWriteReply) :
    replacerWriteString none input mode script = none := rfl

theorem replacerWriteString_initial (pairs : ReplacerPairs) (input : List (BitVec 8)) (mode : Bool)
    (script : List ReplacerWriteReply) :
    replacerWriteString (some (.pending pairs)) input mode script =
      replacerWriteString (some (.ready (replacerBuild pairs))) input mode script := rfl

theorem replacerWriteString_emptyPairs (mode : Bool) (script : List ReplacerWriteReply) :
    (replacerWriteString (some (.pending [])) [] mode script).map Prod.snd =
      some ⟨0, 0, []⟩ := rfl

theorem replacerWriteString_failedWriteKeepsImplementation :
    (match replacerWriteString (some (.pending [([97], [98])])) [97] true [⟨some 0, 1⟩] with
      | some (.ready _, result) => result.error
      | _ => 0) = 1 := rfl

end GoLean.Strings
