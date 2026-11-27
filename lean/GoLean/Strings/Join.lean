import Init.Data.List.Basic
import Init.Data.Nat.Lemmas

namespace GoLean.Strings

private def addJoinLengths (maxInt : Nat) (total : Nat) : List Nat → Option Nat
  | [] => some total
  | size :: sizes =>
    if maxInt - total < size then none
    else addJoinLengths maxInt (total + size) sizes

/-- The size of Go's signed int is represented by maxInt. None is the
    output length overflow panic. The separator guard uses inclusive >=. -/
def joinLength (maxInt : Nat) (sizes : List Nat) (separatorSize : Nat) : Option Nat :=
  match sizes with
  | [] => some 0
  | [one] => some one
  | _ :: _ :: _ =>
    let separators := sizes.length - 1
    if separatorSize != 0 && maxInt / separators ≤ separatorSize then none
    else addJoinLengths maxInt (separatorSize * separators) sizes

private def joinTail (separator : List Nat) : List (List Nat) → List Nat
  | [] => []
  | part :: parts => separator ++ part ++ joinTail separator parts

/-- Go 1.27.0 strings/strings.go:468–496. Size validation precedes the
    Builder appends. Allocation and backing storage are abstracted. -/
def join (maxInt : Nat) (parts : List (List Nat)) (separator : List Nat) : Option (List Nat) :=
  match parts with
  | [] => some []
  | [one] => some one
  | first :: second :: rest => do
    let _ ← joinLength maxInt (parts.map List.length) separator.length
    return first ++ joinTail separator (second :: rest)

theorem join_empty (maxInt : Nat) (separator : List Nat) :
    join maxInt [] separator = some [] := rfl

theorem join_single (maxInt : Nat) (part separator : List Nat) :
    join maxInt [part] separator = some part := rfl

theorem join_two (maxInt : Nat) (first second separator : List Nat) (size : Nat)
    (fits : joinLength maxInt [first.length, second.length] separator.length = some size) :
    join maxInt [first, second] separator = some (first ++ separator ++ second) := by
  simp [join, fits, joinTail, List.append_assoc]

theorem join_separatorOverflow (maxInt : Nat) (first second sep : List Nat)
    (nonempty : sep ≠ []) (limit : maxInt ≤ sep.length) :
    join maxInt [first, second] sep = none := by
  have pos : 0 < sep.length := List.length_pos_iff.mpr nonempty
  simp [join, joinLength, show sep.length ≠ 0 by omega, limit]

theorem join_elementsOverflow (maxInt : Nat) (first second : List Nat)
    (fits : first.length ≤ maxInt) (overflow : maxInt < first.length + second.length) :
    join maxInt [first, second] [] = none := by
  have tooLong : maxInt - first.length < second.length := by omega
  simp [join, joinLength, addJoinLengths, Nat.not_lt.mpr fits, tooLong]

end GoLean.Strings
