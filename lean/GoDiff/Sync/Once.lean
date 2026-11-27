import GoLean.Sync.Once
import GoDiff.Math.Bits.Tuple

open GoLean.Sync

private def fault : OnceFault → String
  | .nilCallback => "nilCallback"
  | .nilPanic => "nilPanic"
  | .payload code => s!"panic:{code}"
  | .nilReceiver => "nilReceiver"

private def result (render : α → String) : Except OnceFault α → String
  | .ok x => render x
  | .error err => fault err

private def callback (mode payload : Int) (value : α) : OnceCallback α :=
  if mode = 0 then .missing else if mode = 2 then .raises (.payload payload)
  else if mode = 3 then .raises .nilPanic else .returns value

private def onceWire (r : OnceObservation) :=
  s!"{result (fun _ => "ok") r.result}:{r.calls}:{r.done.map toString |>.getD "nil"}"

private def twice (s : OnceClosure α) (render : α → String) :=
  let (called, first, next) := callOnceClosure s
  let (again, second, _) := callOnceClosure next
  s!"{result render first}:{called},{result render second}:{again}"

def main : IO Unit := runTuple fun row => do
  let [tag, mode, second, payload, x, y] := row |
    throw (IO.userError "expected receiver, callback modes, panic and results")
  let done := if tag = 0 then none else some (tag == 2)
  let one := callback mode payload ()
  let two := callback second payload ()
  let d := onceDo done one
  let s := onceDoSlow done one
  return String.intercalate "|" [
    onceWire d ++ "," ++ onceWire (onceDo d.done two),
    onceWire s ++ "," ++ onceWire (onceDoSlow s.done two),
    twice (onceFunc one) (fun _ => "ok"),
    twice (onceValue (callback mode payload x)) toString,
    twice (onceValues (callback mode payload (x,y))) (fun (a,b) => s!"{a}/{b}")]
