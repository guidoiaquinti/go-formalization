import GoLean.Strings.Count

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:57–74. The byte slice and string versions
    share the three result branches. The model abstracts Index probes
    and single-byte count acceleration. -/
def count (input separator : List Nat) : Nat :=
  GoLean.Strings.count input separator

theorem count_empty (input : List Nat) :
    count input [] = GoLean.Strings.countRunes input + 1 :=
  GoLean.Strings.count_empty_spec input

theorem count_byte (input : List Nat) (byte : Nat) :
    count input [byte] = input.countP (· == byte) :=
  GoLean.Strings.count_byte_spec input byte

theorem count_nonoverlap (input : List Nat) (first second : Nat) (rest : List Nat) :
    count input (first :: second :: rest) =
      GoLean.Strings.greedyCount (first :: second :: rest) (by simp) input :=
  GoLean.Strings.count_nonoverlap_spec input first second rest

end GoLean.Bytes
