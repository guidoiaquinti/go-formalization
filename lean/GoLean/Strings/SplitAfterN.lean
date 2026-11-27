import GoLean.Strings.SplitAfter

namespace GoLean.Strings

/-- Go 1.27.0 strings/strings.go:326–328. Retain the complete separator. -/
def splitAfterN (input separator : List Nat) (limit : Int) : Option (List (List Nat)) :=
  genSplit input separator separator.length limit

theorem splitAfterN_zero (input separator : List Nat) :
    splitAfterN input separator 0 = none := by
  simp [splitAfterN, genSplit]

theorem splitAfterN_noSeparator (input : List Nat) (limit : Int) :
    splitAfterN input [] limit = splitN input [] limit := rfl

theorem splitAfterN_unbounded (input separator : List Nat) :
    splitAfterN input separator (-1) = splitAfter input separator := rfl

end GoLean.Strings
