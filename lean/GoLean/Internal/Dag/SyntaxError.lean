namespace GoLean.Internal.Dag

inductive SyntaxErrorOutcome where
  | nilReceiver
  | syntax (payload : List UInt8)
  deriving DecidableEq

private def bytes (s : String) : List UInt8 := s.toUTF8.data.toList

/-- Go evaluates the receiver fields before formatting and raising the typed
panic. Model the two string operands as bytes, including malformed UTF-8.
https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L188-L192
https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L237-L247 -/
def rulesSyntaxError (parser : Option (Int × List UInt8)) (msg : List UInt8) :
    SyntaxErrorOutcome :=
  match parser with
  | none => .nilReceiver
  | some (line, lastWord) => .syntax ([
      bytes "parsing graph: line ", bytes (toString line),
      bytes ": syntax error: ", msg, bytes " near ", lastWord].flatten)

theorem rulesSyntaxError_spec (line : Int) (msg lastWord : List UInt8) :
    rulesSyntaxError (some (line, lastWord)) msg =
      .syntax (bytes "parsing graph: line " ++ bytes (toString line) ++
        bytes ": syntax error: " ++ msg ++ bytes " near " ++ lastWord) := by
  simp [rulesSyntaxError, List.append_assoc]

theorem rulesSyntaxError_nil (msg : List UInt8) :
    rulesSyntaxError none msg = .nilReceiver := rfl

end GoLean.Internal.Dag
