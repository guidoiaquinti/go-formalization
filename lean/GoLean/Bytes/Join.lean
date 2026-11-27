import GoLean.Strings.Join

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:551–581. None is the size overflow;
    some none is the nil result when copying an empty singleton. -/
def join (maxInt : Nat) (parts : List (List Nat)) (sep : List Nat) :
    Option (Option (List Nat)) :=
  match parts with
  | [] => some (some [])
  | [part] => some (if part.isEmpty then none else some part)
  | _ :: _ :: _ => (GoLean.Strings.join maxInt parts sep).map some

theorem join_empty (maxInt : Nat) (sep : List Nat) :
    join maxInt [] sep = some (some []) := rfl

theorem join_single (maxInt : Nat) (part sep : List Nat) :
    join maxInt [part] sep = some (if part.isEmpty then none else some part) := rfl

theorem join_two (maxInt : Nat) (first second sep : List Nat) (size : Nat)
    (fits : GoLean.Strings.joinLength maxInt [first.length, second.length] sep.length = some size) :
    join maxInt [first, second] sep = some (some (first ++ sep ++ second)) := by
  simp [join, GoLean.Strings.join_two maxInt first second sep size fits]

end GoLean.Bytes
