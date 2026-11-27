import GoLean.Net.Url.ErrorMethods
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Url

private def inputBytes (row : List Int) : IO (ErrorBytes × List Int) := do
  let length :: rest := row | throw (IO.userError "expected byte count")
  return ((rest.take length.toNat).map Int.toNat, rest.drop length.toNat)

private def condition (code : Int) : Option (Except Nat Bool) :=
  if code == 0 then none else if code == 3 then some (.error 77) else some (.ok (code == 2))

private def encode (bytes : ErrorBytes) := String.intercalate "," (bytes.map toString)

private def render [ToString α] (x : Except ErrorFault α) : String :=
  match x with
  | .ok v => toString v
  | .error .nilReceiver => "nilpanic"
  | .error (.callback token) => s!"callback{token}"

def main : IO Unit := runTuple fun fields => do
  let present :: identity :: timeout :: temporary :: rest := fields |
    throw (IO.userError "expected receiver, identity and capabilities")
  let (op, rest) ← inputBytes rest
  let (url, rest) ← inputBytes rest
  let (inner, rest) ← inputBytes rest
  let (bad, _) ← inputBytes rest
  let receiver : Option ErrorPayload :=
    if present == 0 then none
    else some ⟨identity.toNat, op, url, inner, condition timeout, condition temporary⟩
  return String.intercalate "|" [
    render ((urlErrorError receiver).map encode), render (urlErrorUnwrap receiver),
    render (urlErrorTimeout receiver), render (urlErrorTemporary receiver),
    encode (escapeErrorError bad), encode (invalidHostErrorError bad)]
