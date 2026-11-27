import GoLean.Slices.Order2
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

private def floatKey (word : Nat) : Nat :=
  if word < 2 ^ 63 then word + 2 ^ 63 else 2 ^ 64 - 1 - word

private def floatLess (x y : BitVec 64) : Bool :=
  let a := x.toNat % 2 ^ 63
  let b := y.toNat % 2 ^ 63
  let nan := 0x7ff0000000000000
  if nan < a then decide (b ≤ nan)
  else if nan < b || (a = 0 && b = 0) then false
  else decide (floatKey x.toNat < floatKey y.toNat)

private def wire (result : Order2Observation (BitVec 64) Int) : String :=
  let report := match result.result with
    | .ok (a, b) => s!"ok:{a},{b}"
    | .error .nilCounter => "nil"
    | .error .nilComparator => "nil"
    | .error (.callback payload) => s!"panic:{payload}"
  let calls := String.intercalate ";" (result.calls.map fun (a,b) => s!"{a.toNat},{b.toNat}")
  s!"{report}/{(result.counter.map toString).getD "nil"}/{calls}"

def main : IO Unit := runTuple fun fields => do
  let width :: kind :: a :: b :: present :: counter :: callback :: reply :: raw := fields |
    throw (IO.userError "expected widths, operands, counter, callback and words")
  let words ← raw.mapM (wordInput 64)
  let some left := words[a.toNat]? | throw (IO.userError "invalid left element")
  let some right := words[b.toNat]? | throw (IO.userError "invalid right element")
  let count := if present == 0 then none else some counter
  let less := if kind == 1 then floatLess else fun x y => decide (x.toInt < y.toInt)
  let compare := if callback == 0 then none else some (fun _ _ =>
    if callback == 2 then Except.error reply else Except.ok reply)
  return wire (order2Ordered width.toNat a.toNat b.toNat left right count less) ++ "|" ++
    wire (order2CmpFunc width.toNat a.toNat b.toNat left right count compare)
