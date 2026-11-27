import GoLean.Internal.Dag.SyntaxErrorMethod
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let 0 :: input := fields | throw (IO.userError "expected zero marker and error bytes")
  let bytes ← input.mapM fun x => do return (← wordInput 8 x).toNat.toUInt8
  let result := GoLean.Internal.Dag.syntaxErrorError ⟨bytes⟩
  return String.intercalate " " ("error" :: result.map (fun b => toString b.toNat))
