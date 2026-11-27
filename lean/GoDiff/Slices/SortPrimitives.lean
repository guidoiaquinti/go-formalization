import GoLean.Slices.SortPrimitives
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

def main : IO Unit := runTuple fun fields => do
  let [width, length, hasState, raw] := fields |
    throw (IO.userError "expected width, length, state presence, and word")
  let word ← wordInput 64 raw
  let state := xorshiftNext (if hasState == 0 then none else some word)
  let next := state.map (fun (a, b) => s!"{a.toNat},{b.toNat}") |>.getD "panic"
  let w := word.toNat
  return s!"{next}/{(nextPowerOfTwo width.toNat length).toNat}/{sortIsNaN 23 8 true w}/{sortIsNaN 52 11 true w}/{sortIsNaN 52 11 false w}"
