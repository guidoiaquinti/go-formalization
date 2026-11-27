import GoLean.Strings.ASCIISetContains
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let receiver ← match row with
    | [0] => pure none
    | [1, a, b, c, d] => do
      let words ← [a, b, c, d].mapM (wordInput 64)
      let table := fun (key : BitVec 8) =>
        (words[key.toNat / 64]?.getD 0).toNat.testBit (key.toNat % 64)
      pure (some table)
    | _ => throw (IO.userError "expected nil marker or four table words")
  let entries := (List.range 256).map fun i =>
    match GoLean.Strings.asciiSetContains receiver (BitVec.ofNat 8 i) with
    | none => 'p'
    | some true => '1'
    | some false => '0'
  return String.ofList entries
