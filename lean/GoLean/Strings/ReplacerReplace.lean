import GoLean.Strings.ReplacerBuildOnce
import GoLean.Strings.SingleStringReplacerReplace
import GoLean.Strings.GenericReplacerReplace
import GoLean.Strings.ByteReplacerReplace

namespace GoLean.Strings

inductive ReplacerState where
  | pending (pairs : ReplacerPairs)
  | ready (implementation : BuiltReplacer)

/-- Sequential view of the first call and later calls using sync.Once. -/
def replacerImplementation : ReplacerState → BuiltReplacer
  | .pending pairs => replacerBuild pairs
  | .ready implementation => implementation

private def replaceBuilt (implementation : BuiltReplacer) (input : List (BitVec 8)) :
    Option (List (BitVec 8)) :=
  match implementation with
  | .single state => some (singleStringReplacerReplace state.finder.pattern state.value input)
  | .generic pairs _ => some (genericReplacerReplace pairs input)
  | .bytes mapping => byteReplacerReplace (some mapping) input
  | .byteStrings state => byteStringReplacerReplace (some state) input

/-- Go 1.27.0 src/strings/replace.go:94–98. Build the implementation
    on first use and dispatch to it. Later calls keep the implementation. -/
def replacerReplace (receiver : Option ReplacerState) (input : List (BitVec 8)) :
    Option (ReplacerState × List (BitVec 8)) := do
  let state ← receiver
  let implementation := replacerImplementation state
  let result ← replaceBuilt implementation input
  return (.ready implementation, result)

theorem replacerReplace_nil (input : List (BitVec 8)) : replacerReplace none input = none := rfl

theorem replacerReplace_initial (pairs : ReplacerPairs) (input : List (BitVec 8)) :
    replacerReplace (some (.pending pairs)) input =
      replacerReplace (some (.ready (replacerBuild pairs))) input := rfl

theorem replacerReplace_emptyPairs (input : List (BitVec 8)) :
    (replacerReplace (some (.pending [])) input).map Prod.snd = some input := by
  simp [replacerReplace, replacerImplementation, replacerBuild_empty, replaceBuilt,
    byteReplacerReplace_identity]

theorem replacerReplace_priority :
    (replacerReplace (some (.pending [([97], [98]), ([97], [99])])) [97]).map Prod.snd =
      some [98] := rfl

end GoLean.Strings
