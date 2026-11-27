import GoLean.Sort.BreakPatterns
import GoDiff.Math.Bits.Tuple

open GoLean.Sort

def main : IO Unit := runTuple fun row => do
  let [kind, width, a, b] := row |
    throw (IO.userError "expected helper, word width and arguments")
  let width := width.toNat
  match kind with
  | 0 =>
    let state ← wordInput 64 a
    return toString (xorshiftNext state).toNat
  | 1 => return toString (nextPowerOfTwo (BitVec.ofInt width a)).toNat
  | _ =>
    let trace := if kind = 2 then breakPatterns width a.toNat b.toNat
      else breakPatternsFunc width a.toNat b.toNat
    return "swaps:" ++ String.intercalate ";" (trace.map fun (i, j) => s!"{i},{j}")
