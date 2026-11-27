import GoLean.Internal.Dag.Parse
import GoDiff.Math.Bits.Tuple

open GoLean.Internal.Dag

private def encode (value : List UInt8) : String :=
  String.intercalate " " (toString value.length :: value.map (fun b => toString b.toNat))

def main : IO Unit := runTuple fun fields => do
  let width :: words := fields | throw (IO.userError "expected native int width and input bytes")
  let input ← words.mapM fun n => do
    return (← wordInput 8 n).toNat.toUInt8
  return match parse width.toNat input with
  | .exhausted => "exhausted"
  | .syntax (.syntax payload) => "syntax " ++ encode payload
  | .syntax .nilReceiver => "nil receiver"
  | .semantic err => "semantic " ++ encode err
  | .graph nodes arcs =>
    let rows := nodes.flatMap fun source =>
      nodes.map fun target => if arcs.contains (source, target) then "1" else "0"
    String.intercalate " " ("graph" :: toString nodes.length :: nodes.map encode ++ rows)
