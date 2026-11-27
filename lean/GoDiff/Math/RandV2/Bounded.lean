import GoLean.Math.RandV2.Bounded
import GoDiff.Math.Bits.Tuple

open GoLean.Math.RandV2

private def wire (r : BoundedResult) := s!"{r.status}:{r.value}:{r.pulls}"

def main : IO Unit := runTuple fun row => do
  let width :: rawBound :: raw := row | throw (IO.userError "expected word width, bound and source words")
  let bound := (← wordInput 64 rawBound).toNat
  let words ← raw.mapM fun w => BitVec.toNat <$> wordInput 64 w
  let width := width.toNat
  return String.intercalate "|" ([
    randUint64n width bound words, randUint32n (bound % 2^32) words,
    randInt64N bound words, randUint64N bound words, randInt32N bound words,
    randUint32N bound words, randIntN width bound words, randUintN width bound words,
    randN true 8 bound words, randN false 16 bound words].map wire)
