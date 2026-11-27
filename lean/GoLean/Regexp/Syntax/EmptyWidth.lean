import Init.Data.BitVec.Basic

namespace GoLean.Regexp.Syntax

/-- Go 1.27.0 src/regexp/syntax/prog.go:108–112. Word boundaries use ASCII. -/
def isWordChar (r : Int) : Bool :=
  (97 ≤ r && r ≤ 122) || (65 ≤ r && r ≤ 90) || (48 ≤ r && r ≤ 57) || r == 95

/-- Go 1.27.0 prog.go:80–103. Any negative rune marks a text boundary here. -/
def emptyOpContext (before after : Int) : BitVec 8 :=
  let left : BitVec 8 := if isWordChar before then 0 else
    if before = 10 then 1 else if before < 0 then 5 else 0
  let right : BitVec 8 := if isWordChar after then 0 else
    if after = 10 then 2 else if after < 0 then 10 else 0
  let boundary : BitVec 8 := if isWordChar before != isWordChar after then 16 else 32
  left ||| right ||| boundary

inductive EmptyWidthOutcome where
  | matched (answer : Bool)
  | nilReceiver
  | unknownArg
  deriving DecidableEq

/-- Go 1.27.0 prog.go:266–282. Match one known low-byte operation.
    This entry recognizes -1 as the text boundary and ignores the opcode field. -/
def instMatchEmptyWidth (argument : Option (BitVec 32)) (before after : Int) :
    EmptyWidthOutcome :=
  match argument with
  | none => .nilReceiver
  | some arg =>
    match arg.toNat % 256 with
    | 1 => .matched (before == 10 || before == -1)
    | 2 => .matched (after == 10 || after == -1)
    | 4 => .matched (before == -1)
    | 8 => .matched (after == -1)
    | 16 => .matched (isWordChar before != isWordChar after)
    | 32 => .matched (isWordChar before == isWordChar after)
    | _ => .unknownArg

/-- Go 1.27.0 prog.go:138–146. Collapse the three specialized rune opcodes. -/
def instOp (receiver : Option (BitVec 8)) : Option (BitVec 8) :=
  receiver.map fun op => if op == 8 || op == 9 || op == 10 then 7 else op

theorem isWordChar_underscore : isWordChar 95 = true := rfl
theorem isWordChar_nonASCII : isWordChar 0x212a = false := rfl

theorem emptyOpContext_negative :
    emptyOpContext (-9) (-2) = 47 := rfl

theorem emptyOpContext_newlineWord :
    emptyOpContext 10 95 = 17 := rfl

theorem instMatchEmptyWidth_exactSentinel :
    instMatchEmptyWidth (some 260) (-2) (-1) = .matched false ∧
    instMatchEmptyWidth (some 260) (-1) (-2) = .matched true := ⟨rfl, rfl⟩

theorem instMatchEmptyWidth_mixture :
    instMatchEmptyWidth (some 3) (-1) (-1) = .unknownArg := rfl

theorem instOp_runes :
    instOp (some 8) = some 7 ∧ instOp (some 9) = some 7 ∧ instOp (some 10) = some 7 :=
  ⟨rfl, rfl, rfl⟩

theorem instOp_unknown : instOp (some 255) = some 255 := rfl

end GoLean.Regexp.Syntax
