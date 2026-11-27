import GoLean.Regexp.Syntax.OpcodeStrings
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [raw] := fields | throw (IO.userError "expected a byte opcode")
  let op ← wordInput 8 raw
  return GoLean.Regexp.Syntax.instOpString op ++ "|" ++ GoLean.Regexp.Syntax.opString op
