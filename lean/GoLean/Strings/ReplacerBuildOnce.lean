import GoLean.Strings.ReplacerBuild

namespace GoLean.Strings

structure ReplacerBuildOnceResult where
  pending : Option ReplacerPairs
  onceIsZero : Bool
  implementation : BuiltReplacer

/-- Go 1.27.0 src/strings/replace.go:39–42. The callback builds from
    pending arguments and clears them. Calling it directly leaves sync.Once
    unchanged. None as the receiver models the nil-receiver panic. -/
def replacerBuildOnce (receiver : Option (Option ReplacerPairs × Bool)) :
    Option ReplacerBuildOnceResult :=
  receiver.map fun (pending, onceIsZero) =>
    ⟨none, onceIsZero, replacerBuild (pending.getD [])⟩

theorem replacerBuildOnce_nil : replacerBuildOnce none = none := rfl

theorem replacerBuildOnce_pending (pairs : Option ReplacerPairs) (onceIsZero : Bool) :
    (replacerBuildOnce (some (pairs, onceIsZero))).map ReplacerBuildOnceResult.pending =
      some none := rfl

theorem replacerBuildOnce_once (pairs : Option ReplacerPairs) (onceIsZero : Bool) :
    (replacerBuildOnce (some (pairs, onceIsZero))).map ReplacerBuildOnceResult.onceIsZero =
      some onceIsZero := rfl

theorem replacerBuildOnce_implementation (pairs : ReplacerPairs) (onceIsZero : Bool) :
    (replacerBuildOnce (some (some pairs, onceIsZero))).map ReplacerBuildOnceResult.implementation =
      some (replacerBuild pairs) := rfl

theorem replacerBuildOnce_noPending (onceIsZero : Bool) :
    (replacerBuildOnce (some (none, onceIsZero))).map ReplacerBuildOnceResult.implementation =
      some (.bytes id) := rfl

end GoLean.Strings
