import GoLean.Image.Color.ModelFunc
import GoLean.Image.Color.ModelFuncConvert
import GoDiff.Math.Bits.Tuple

open GoLean.Image.Color

private def observe (m : ModelFuncWrapper (Nat → Nat → Nat × Except Nat Nat))
    (arg : Nat) : String :=
  let (calls, _, result) := modelFuncConvert (some m.callback) arg 0
  match result with
  | .ok value => s!"ok {value} {calls.length} {calls.headD 0}"
  | .error .nilCallback => "nil-panic 0 0 0"
  | .error .nilReceiver => "nil-receiver 0 0 0"
  | .error (.callback value) => s!"callback-panic {value} {calls.length} {calls.headD 0}"

def main : IO Unit := runTuple fun fields => do
  let [kind, arg, capture] := fields | throw (IO.userError "expected kind, argument, capture")
  let callback (input state : Nat) : Nat × Except Nat Nat :=
    (state + 1, if kind = 1 then .ok (input + capture.toNat) else .error capture.toNat)
  let callback := if kind = 0 then none else some callback
  let (first, next) := modelFunc 0 callback
  let (second, _) := modelFunc next callback
  return s!"pointer nonnil {if first.identity = second.identity then "same" else "distinct"} 0 {observe first arg.toNat} {observe second arg.toNat}"
