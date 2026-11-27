import GoLean.Regexp.Syntax.EmptyWidth
import GoDiff.Math.Bits.Tuple

open GoLean.Regexp.Syntax

def main : IO Unit := runTuple fun fields => do
  let [before, after, present, op, arg] := fields |
    throw (IO.userError "expected two runes and instruction fields")
  let before := (BitVec.ofInt 32 before).toInt
  let after := (BitVec.ofInt 32 after).toInt
  let arg := if present == 0 then none else some (BitVec.ofInt 32 arg)
  let op := if present == 0 then none else some (BitVec.ofInt 8 op)
  let matched := match instMatchEmptyWidth arg before after with
    | .matched ok => toString ok
    | .nilReceiver => "nil"
    | .unknownArg => "unknown"
  let normalized := (instOp op).map (toString ∘ BitVec.toNat) |>.getD "nil"
  return s!"{isWordChar before}|{isWordChar after}|{(emptyOpContext before after).toNat}|{matched}|{normalized}"
