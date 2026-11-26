import GoLean.Math.Bits.Division
import GoDiff.Math.Bits.Tuple

def main : IO Unit :=
  runTuple fun fields => do
    let [hi, lo, y] := fields | throw (IO.userError "wrong number of inputs")
    let hi ← wordInput 32 hi
    let lo ← wordInput 32 lo
    let y ← wordInput 32 y
    match GoLean.Math.Bits.div32 hi lo y with
    | .error .divideByZero => return "divide-by-zero"
    | .error .quotientOverflow => return "quotient-overflow"
    | .ok (q, r) => return s!"{q.toNat} {r.toNat}"
