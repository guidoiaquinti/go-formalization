import GoLean.Internal.Dag.SyntaxError
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let flag :: line :: msgLength :: input := fields |
    throw (IO.userError "expected receiver flag, line, message length, and bytes")
  if msgLength < 0 || input.length < msgLength.toNat then
    throw (IO.userError "invalid message length")
  let bytes ← input.mapM fun x => do return (← wordInput 8 x).toNat.toUInt8
  let parser := if flag == 0 then none else some (line, bytes.drop msgLength.toNat)
  match GoLean.Internal.Dag.rulesSyntaxError parser (bytes.take msgLength.toNat) with
  | .nilReceiver => return "nil"
  | .syntax payload =>
    return String.intercalate " " ("syntax" :: payload.map (fun b => toString b.toNat))
