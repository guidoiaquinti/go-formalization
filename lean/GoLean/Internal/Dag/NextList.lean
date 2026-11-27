import GoLean.Internal.Dag.NextToken

namespace GoLean.Internal.Dag

inductive ListOutcome (width : Nat) where
  | returned (names : List (List UInt8)) (token : List UInt8) (state : TokenState width)
  | panicked (reason : SyntaxErrorOutcome) (state : Option (TokenState width))
  | exhausted

private def separator (tok : List UInt8) : Bool :=
  [[44], [60], [33, 60], [59]].contains tok

private def badList (state : TokenState width) : ListOutcome width :=
  .panicked (rulesSyntaxError (some (state.line.toInt, state.lastWord))
    ("bad list syntax".toUTF8.data.toList)) (some state)

/-- Go accumulates a name and reads one following token, looping only after
a comma. Each continuation consumes at least that comma from the input.
An empty names list represents Go's nil EOF result.
https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L249-L269 -/
private def nextListLoop (width : Nat) :
    Nat → List (List UInt8) → Option (TokenState width) → ListOutcome width
  | 0, _, _ => .exhausted
  | fuel + 1, names, src =>
    match nextToken src with
    | .panicked reason state => .panicked reason state
    | .returned name state =>
      if name.isEmpty && names.isEmpty then .returned [] [] state
      else if name.isEmpty || separator name then badList state
      else
        let names := name :: names
        match nextToken (some state) with
        | .panicked reason state => .panicked reason state
        | .returned tok state =>
          if tok == [44] then nextListLoop width fuel names (some state)
          else .returned names.reverse tok state

def nextList (src : Option (TokenState width)) : ListOutcome width :=
  nextListLoop width ((src.map (·.text.length)).getD 0 + 1) [] src

theorem nextList_nil (width : Nat) :
    nextList (width := width) none = .panicked .nilReceiver none := rfl

theorem nextList_eof (src state : TokenState width)
    (read : nextToken (some src) = .returned [] state) :
    nextList (some src) = .returned [] [] state := by
  simp [nextList, nextListLoop, read]

theorem nextList_first_error (src : TokenState width) (reason : SyntaxErrorOutcome)
    (state : Option (TokenState width))
    (read : nextToken (some src) = .panicked reason state) :
    nextList (some src) = .panicked reason state := by
  simp [nextList, nextListLoop, read]

theorem nextList_single (src afterName final : TokenState width) (name tok : List UInt8)
    (nonempty : name ≠ []) (valid : name ∉ [[44], [60], [33, 60], [59]])
    (done : tok ≠ [44])
    (first : nextToken (some src) = .returned name afterName)
    (second : nextToken (some afterName) = .returned tok final) :
    nextList (some src) = .returned [name] tok final := by
  simp [nextList, nextListLoop, first, second, separator, nonempty, valid, done]

/-- A delimiter in the first name position raises the typed list error
after the lexer has already updated the receiver. -/
theorem nextList_badDelimiter (src state : TokenState width) (tok : List UInt8)
    (invalid : tok ∈ [[44], [60], [33, 60], [59]])
    (read : nextToken (some src) = .returned tok state) :
    nextList (some src) = .panicked
      (rulesSyntaxError (some (state.line.toInt, state.lastWord))
        ("bad list syntax".toUTF8.data.toList)) (some state) := by
  simp at invalid
  rcases invalid with rfl | rfl | rfl | rfl <;>
    simp [nextList, nextListLoop, read, separator, badList]

/-- Once a comma is consumed, EOF is an error even if the lexer has reached
a later line and changed the remaining text. -/
theorem nextList_trailingComma (src afterName afterComma final : TokenState width) (name : List UInt8)
    (size : src.text.length ≠ 0) (nonempty : name ≠ [])
    (valid : name ∉ [[44], [60], [33, 60], [59]])
    (first : nextToken (some src) = .returned name afterName)
    (second : nextToken (some afterName) = .returned [44] afterComma)
    (third : nextToken (some afterComma) = .returned [] final) :
    nextList (some src) = .panicked
      (rulesSyntaxError (some (final.line.toInt, final.lastWord))
        ("bad list syntax".toUTF8.data.toList)) (some final) := by
  obtain ⟨fuel, h⟩ := Nat.exists_eq_succ_of_ne_zero size
  simp [nextList, nextListLoop, h, first, second, third, separator, nonempty, valid, badList]

end GoLean.Internal.Dag
