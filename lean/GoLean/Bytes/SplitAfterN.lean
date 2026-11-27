import GoLean.Bytes.SplitN
import GoLean.Strings.SplitAfterN

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:419–421 retains the full separator.
    It caps processed fragments but preserves capacity for a raw suffix. -/
def splitAfterN (input separator : List Nat) (spare : Nat) (limit : Int) :
    Option (List (List Nat × Nat)) :=
  withSplitCaps input separator spare limit (GoLean.Strings.splitAfterN input separator limit)

theorem splitAfterN_zero (input separator : List Nat) (spare : Nat) :
    splitAfterN input separator spare 0 = none := by
  simp [splitAfterN, withSplitCaps, GoLean.Strings.splitAfterN_zero]

theorem splitAfterN_noSeparator (input : List Nat) (spare : Nat) (limit : Int) :
    splitAfterN input [] spare limit = splitN input [] spare limit := by
  simp [splitAfterN, splitN, GoLean.Strings.splitAfterN_noSeparator]

end GoLean.Bytes
