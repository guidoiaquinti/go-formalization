import GoLean.Strings.GetStringWriter
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let mode :: byteLimit :: byteError :: stringLimit :: stringError :: fields := fields
    | throw (IO.userError "expected selection mode, method replies and payload")
  let payload ← fields.mapM (wordInput 8)
  let respond := fun cap error (raw : List (BitVec 8)) =>
    (min cap (Int.ofNat raw.length), error.toNat)
  let byteMethod := respond byteLimit byteError
  let stringMethod := respond stringLimit stringError
  let capability := if mode == 1 || mode == 2 then some stringMethod else none
  let selected := GoLean.Strings.getStringWriter byteMethod capability
  let route := if selected.original then "1" else "0"
  if mode > 1 then
    return route
  let (count, error) := selected.writeString payload
  let event := (if selected.original then "S" else "B") ++
    String.join (payload.map (fun byte => "," ++ toString byte.toNat))
  return s!"{route}:{count}:{error}:{event}"
