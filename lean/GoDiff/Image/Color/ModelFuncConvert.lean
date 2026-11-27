import GoLean.Image.Color.ModelFuncConvert
import GoDiff.Math.Bits.Tuple

open GoLean.Image.Color

def main : IO Unit := runTuple fun fields => do
  let [kind, input, initial, delta] := fields |
    throw (IO.userError "expected mode, input token, state and payload")
  let c := input.toNat
  let callback (arg state : Nat) : Nat × Except Nat Nat :=
    (state + delta.toNat + arg,
      if kind = 2 then .ok arg
      else if kind = 3 then .ok 0
      else if kind = 4 then .ok 4294967297
      else .error delta.toNat)
  let m := if kind = 0 then none else if kind = 1 then some none else some (some callback)
  let (calls, state, result) := modelFuncConvert m c initial.toNat
  let (tag, payload) :=
    match result with
    | .ok value => ("ok", value)
    | .error .nilReceiver | .error .nilCallback => ("nil-panic", 0)
    | .error (.callback value) => ("callback-panic", value)
  return s!"{tag} {payload} {calls.length} {calls.headD 0} {state}"
