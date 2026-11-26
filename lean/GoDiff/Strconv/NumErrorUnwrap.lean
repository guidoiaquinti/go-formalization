import GoLean.Strconv.NumErrorUnwrap
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [tag, func, num] := fields | throw (IO.userError "expected error tag and two metadata bytes")
  let receiver := if tag = 0 then none else some
    (⟨[func.toNat.toUInt8], [num.toNat.toUInt8], tag⟩ : GoLean.Strconv.NumError Int)
  match GoLean.Strconv.numErrorUnwrap receiver with
  | .error () => return "nil receiver panic"
  | .ok value => return toString value
