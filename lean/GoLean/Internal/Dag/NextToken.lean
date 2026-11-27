import GoLean.Internal.Dag.SyntaxError
import Init.Data.BitVec.Lemmas

namespace GoLean.Internal.Dag

structure TokenState (width : Nat) where
  line : BitVec width
  lastWord : List UInt8
  text : List UInt8
  deriving DecidableEq

inductive TokenOutcome (width : Nat) where
  | returned (token : List UInt8) (state : TokenState width)
  | panicked (reason : SyntaxErrorOutcome) (state : Option (TokenState width))
  deriving DecidableEq

private def delimiter (b : UInt8) : Bool :=
  [33, 59, 44, 60, 35, 10, 32, 9].contains b

/-- Skip spaces, tabs, newlines, and comments. A comment leaves its terminating
newline in the suffix for the ordinary newline branch to count. Fuel is at
least the number of bytes, and each continuation consumes at least one. -/
def tokenTrivia (fuel : Nat) (text : List UInt8) : Nat × List UInt8 :=
  match fuel, text with
  | 0, _ => (0, text)
  | _ + 1, [] => (0, [])
  | n + 1, b :: rest =>
    if b == 35 then tokenTrivia n (rest.dropWhile (· != 10))
    else if b == 10 then
      let (lines, suffix) := tokenTrivia n rest
      (lines + 1, suffix)
    else if b == 32 || b == 9 then tokenTrivia n rest
    else (0, text)

/-- Model the next-token scan and its receiver writes. The integer field wraps
at the native Go int width. The old word survives punctuation and panics.
https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L237-L247
https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L271-L317 -/
def nextToken (src : Option (TokenState width)) : TokenOutcome width :=
  match src with
  | none => .panicked .nilReceiver none
  | some input =>
    let (lines, text) := tokenTrivia input.text.length input.text
    let state := { input with line := input.line + BitVec.ofNat width lines, text }
    match text with
    | [] => .returned [] state
    | b :: rest =>
      if b == 59 || b == 44 || b == 60 then
        .returned [b] { state with text := rest }
      else if b == 33 then
        match rest with
        | c :: tail =>
          if c == 60 then .returned [33, 60] { state with text := tail }
          else .panicked (rulesSyntaxError (some (state.line.toInt, state.lastWord))
            ("unexpected token !".toUTF8.data.toList)) (some state)
        | [] => .panicked (rulesSyntaxError (some (state.line.toInt, state.lastWord))
            ("unexpected token !".toUTF8.data.toList)) (some state)
      else
        let word := text.takeWhile (! delimiter ·)
        .returned word { state with text := text.drop word.length, lastWord := word }

private theorem dropWord_eq (text : List UInt8) :
    text.drop (text.takeWhile (! delimiter ·)).length = text.dropWhile (! delimiter ·) := by
  have h := List.takeWhile_append_dropWhile (p := (! delimiter ·)) (l := text)
  have eq := congrArg (List.drop (text.takeWhile (! delimiter ·)).length) h
  simpa only [List.drop_left] using eq.symm

/-- After the trivia transition, a name splits at the first delimiter. The
delimiter stays in the buffer and the full returned word becomes lastWord. -/
theorem nextToken_word (src : TokenState width) (lines : Nat) (b : UInt8) (rest : List UInt8)
    (h : tokenTrivia src.text.length src.text = (lines, b :: rest))
    (hb : delimiter b = false) :
    nextToken (some src) =
      .returned ((b :: rest).takeWhile (! delimiter ·))
        ⟨src.line + BitVec.ofNat width lines, (b :: rest).takeWhile (! delimiter ·),
          (b :: rest).dropWhile (! delimiter ·)⟩ := by
  have hclass : (b == 59 || b == 44 || b == 60) = false ∧ (b == 33) = false := by
    simp [delimiter] at hb
    simp_all
  simp only [nextToken, h, hclass.1, hclass.2, Bool.false_eq_true, ↓reduceIte]
  rw [dropWord_eq]

/-- A malformed bang keeps the offending suffix and previous word while
preserving the line increment already performed by the trivia scan. -/
theorem nextToken_bang_error (src : TokenState width) (lines : Nat) (rest : List UInt8)
    (h : tokenTrivia src.text.length src.text = (lines, 33 :: rest))
    (hb : rest.head? ≠ some 60) :
    let line := src.line + BitVec.ofNat width lines
    nextToken (some src) =
      .panicked (rulesSyntaxError (some (line.toInt, src.lastWord))
        ("unexpected token !".toUTF8.data.toList)) (some ⟨line, src.lastWord, 33 :: rest⟩) := by
  cases rest with
  | nil => simp [nextToken, h]
  | cons c tail =>
    have hc : c ≠ 60 := by simpa using hb
    simp [nextToken, h, hc]

theorem nextToken_eof (src : TokenState width) (lines : Nat)
    (h : tokenTrivia src.text.length src.text = (lines, [])) :
    nextToken (some src) =
      .returned [] ⟨src.line + BitVec.ofNat width lines, src.lastWord, []⟩ := by
  simp [nextToken, h]

theorem nextToken_punctuation (src : TokenState width) (lines : Nat) (b : UInt8)
    (rest : List UInt8) (h : tokenTrivia src.text.length src.text = (lines, b :: rest))
    (hb : b = 59 ∨ b = 44 ∨ b = 60) :
    nextToken (some src) =
      .returned [b] ⟨src.line + BitVec.ofNat width lines, src.lastWord, rest⟩ := by
  rcases hb with rfl | rfl | rfl <;> simp [nextToken, h]

theorem nextToken_bang (src : TokenState width) (lines : Nat) (rest : List UInt8)
    (h : tokenTrivia src.text.length src.text = (lines, 33 :: 60 :: rest)) :
    nextToken (some src) =
      .returned [33, 60] ⟨src.line + BitVec.ofNat width lines, src.lastWord, rest⟩ := by
  simp [nextToken, h]

theorem nextToken_nil (width : Nat) :
    nextToken (width := width) none = .panicked .nilReceiver none := rfl

end GoLean.Internal.Dag
