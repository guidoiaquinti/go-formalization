import GoLean.Time.TimeCompare
import GoDiff.Math.Bits.Tuple

open GoLean.Time

def main : IO Unit := runTuple fun fields => do
  let [tw, te, uw, ue] := fields | throw (IO.userError "expected two clock words per operand")
  let t : ClockWords := ⟨← wordInput 64 tw, BitVec.ofInt 64 te⟩
  let u : ClockWords := ⟨← wordInput 64 uw, BitVec.ofInt 64 ue⟩
  return s!"{clockNsec t} {clockSec t} {clockUnixSec t} {clockMono t} {clockNanosecond t} {clockUnix t} {clockIsZero t} {clockBefore t u} {clockAfter t u} {clockEqual t u} {clockCompare t u}"
