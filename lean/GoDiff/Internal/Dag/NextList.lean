import GoLean.Internal.Dag.NextList
import GoDiff.Math.Bits.Tuple

open GoLean.Internal.Dag

private def encode (input : List UInt8) : String :=
  String.intercalate " " (toString input.length :: input.map (fun b => toString b.toNat))

private def formatState (p : TokenState width) : String :=
  s!"{p.line.toInt} {encode p.lastWord} {encode p.text}"

def main : IO Unit := runTuple fun fields => do
  let width :: flag :: line :: n :: input := fields |
    throw (IO.userError "expected width, receiver flag, line, last-word length, and bytes")
  if n < 0 || input.length < n.toNat then throw (IO.userError "invalid word length")
  let bytes ← input.mapM fun x => do return (← wordInput 8 x).toNat.toUInt8
  let parser : Option (TokenState width.toNat) :=
    if flag == 0 then none else some ⟨BitVec.ofInt _ line, bytes.take n.toNat, bytes.drop n.toNat⟩
  match nextList parser with
  | .returned names tok state =>
    let encoded := String.intercalate " " (toString names.length :: names.map encode)
    return s!"return {encoded} {encode tok} {formatState state}"
  | .panicked .nilReceiver _ => return "nil"
  | .panicked (.syntax payload) (some state) =>
    return s!"syntax {encode payload} {formatState state}"
  | .panicked (.syntax _) none => throw (IO.userError "missing parser state")
  | .exhausted => throw (IO.userError "parser exhausted its byte bound")
