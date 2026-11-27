import GoLean.Internal.Dag.NextList

namespace GoLean.Internal.Dag

structure Rule where
  less : List (List UInt8)
  op : List UInt8
  defn : List (List UInt8)
  deriving DecidableEq

inductive ParseRulesOutcome where
  | returned (rules : List Rule)
  | failed (rules : List Rule) (error : SyntaxErrorOutcome)
  | exhausted

private def failRules (out : List Rule) (state : TokenState width) (message : String) :
    ParseRulesOutcome :=
  .failed out.reverse (rulesSyntaxError (some (state.line.toInt, state.lastWord))
    message.toUTF8.data.toList)

/-- The outer loop starts a new chain after a semicolon. Each continued iteration
consumes a terminator, so input length plus one bounds the number of list reads.
The reversed accumulator preserves the rules already emitted before a syntax error.
https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L181-L235 -/
def parseRulesLoop (width : Nat) :
    Nat → List Rule → Option (List (List UInt8) × List UInt8) →
    TokenState width → ParseRulesOutcome
  | 0, _, _, _ => .exhausted
  | fuel + 1, out, previous, state =>
    match nextList (some state) with
    | .exhausted => .exhausted
    | .panicked reason _ => .failed out.reverse reason
    | .returned names token next =>
      if token.isEmpty then
        if previous.isNone then .returned out.reverse
        else failRules out next "unexpected EOF"
      else
        let updated := match previous with
          | none => out
          | some (less, op) => ⟨less, op, names⟩ :: out
        if token == [59] then parseRulesLoop width fuel updated none next
        else if token == [60] || token == [33, 60] then
          parseRulesLoop width fuel updated (some (names, token)) next
        else failRules updated next "missing <"

def parseRules (width : Nat) (text : List UInt8) : ParseRulesOutcome :=
  parseRulesLoop width (text.length + 1) [] none ⟨1#width, [], text⟩

/-- A closed two-list rule is emitted in order before EOF. -/
theorem parseRules_single (width : Nat) (text : List UInt8) (left right : List (List UInt8))
    (first second final : TokenState width) (op : List UInt8)
    (size : 2 ≤ text.length) (valid : op = [60] ∨ op = [33, 60])
    (readFirst : nextList (some (⟨1#width, [], text⟩ : TokenState width)) = .returned left op first)
    (readSecond : nextList (some first) = .returned right [59] second)
    (readEOF : nextList (some second) = .returned [] [] final) :
    parseRules width text = .returned [⟨left, op, right⟩] := by
  obtain ⟨fuel, amount⟩ : ∃ fuel, text.length = fuel + 2 := ⟨text.length - 2, by omega⟩
  rcases valid with rfl | rfl <;>
    simp [parseRules, amount, parseRulesLoop, readFirst, readSecond, readEOF]

theorem parseRules_listError (width fuel : Nat) (out : List Rule)
    (previous : Option (List (List UInt8) × List UInt8)) (state : TokenState width)
    (reason : SyntaxErrorOutcome) (final : Option (TokenState width))
    (read : nextList (some state) = .panicked reason final) :
    parseRulesLoop width (fuel + 1) out previous state = .failed out.reverse reason := by
  simp [parseRulesLoop, read]

/-- EOF with an unclosed chain reports an error before appending another rule. -/
theorem parseRules_openEOF (width fuel : Nat) (out : List Rule) (left : List (List UInt8))
    (op : List UInt8) (state final : TokenState width) (names : List (List UInt8))
    (read : nextList (some state) = .returned names [] final) :
    parseRulesLoop width (fuel + 1) out (some (left, op)) state =
      .failed out.reverse (rulesSyntaxError (some (final.line.toInt, final.lastWord))
        ("unexpected EOF".toUTF8.data.toList)) := by
  simp [parseRulesLoop, read, failRules]

/-- A second list followed by a name in the operator slot is emitted before
the missing-operator error. A lexical or list error takes the earlier path. -/
theorem parseRules_missingOperator (width : Nat) (text : List UInt8)
    (left right : List (List UInt8)) (first final : TokenState width) (tok : List UInt8)
    (size : 1 ≤ text.length) (invalid : tok ≠ [] ∧ tok ≠ [59] ∧ tok ≠ [60] ∧ tok ≠ [33, 60])
    (readFirst : nextList (some (⟨1#width, [], text⟩ : TokenState width)) = .returned left [60] first)
    (readSecond : nextList (some first) = .returned right tok final) :
    parseRules width text = .failed [⟨left, [60], right⟩]
      (rulesSyntaxError (some (final.line.toInt, final.lastWord))
        ("missing <".toUTF8.data.toList)) := by
  obtain ⟨fuel, amount⟩ : ∃ fuel, text.length = fuel + 1 := ⟨text.length - 1, by omega⟩
  simp [parseRules, amount, parseRulesLoop, readFirst, readSecond,
    invalid.1, invalid.2.1, invalid.2.2.1, invalid.2.2.2, failRules]

end GoLean.Internal.Dag
