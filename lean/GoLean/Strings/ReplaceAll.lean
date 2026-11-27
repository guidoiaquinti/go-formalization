import GoLean.Strings.Replace

namespace GoLean.Strings

/-- Go 1.27.0 src/strings/strings.go:1173–1180. Delegate with a negative
    count so every nonoverlapping match is replaced. -/
def replaceAll (input old replacement : List Nat) : List Nat :=
  replace input old replacement (-1)

theorem replaceAll_unlimited (input old replacement : List Nat) :
    replaceAll input old replacement = replace input old replacement (-1) := rfl

theorem replaceAll_missing (input old replacement : List Nat) (missing : count input old = 0) :
    replaceAll input old replacement = input :=
  replace_missing input old replacement (-1) missing

theorem replaceAll_emptyMalformed :
    replaceAll [0xe2, 0x82, 0xef, 0xbf, 0xbd] [] [0xff] =
      [0xff, 0xe2, 0xff, 0x82, 0xff, 0xef, 0xbf, 0xbd, 0xff] :=
  replace_emptyMalformed

end GoLean.Strings
