import GoLean.Internal.Dag.ParseRules
import GoDiff.Math.Bits.Tuple

open GoLean.Internal.Dag

private def encode (input : List UInt8) : List String :=
  toString input.length :: input.map (fun b => toString b.toNat)

private def encodeList (names : List (List UInt8)) : List String :=
  toString names.length :: names.flatMap encode

private def render (rules : List Rule) : List String :=
  toString rules.length :: rules.flatMap (fun r => encodeList r.less ++ encode r.op ++ encodeList r.defn)

def main : IO Unit := runTuple fun fields => do
  let width :: input := fields | throw (IO.userError "expected width and input bytes")
  let bytes ← input.mapM fun x => do return (← wordInput 8 x).toNat.toUInt8
  let output ← match parseRules width.toNat bytes with
    | .returned rules => pure ("ok" :: render rules)
    | .failed rules (.syntax error) => pure ("error" :: render rules ++ encode error)
    | .failed _ .nilReceiver => throw (IO.userError "unexpected nil receiver")
    | .exhausted => throw (IO.userError "parser exhausted its byte bound")
  return String.intercalate " " output
