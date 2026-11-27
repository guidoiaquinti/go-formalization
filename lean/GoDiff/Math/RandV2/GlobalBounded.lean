import GoLean.Math.RandV2.GlobalBounded
import GoDiff.Math.Bits.Tuple

open GoLean.Math.RandV2

private def wire (r : BoundedResult) := s!"{r.status}:{r.value}:{r.pulls}"

def main : IO Unit := runTuple fun row => do
  let width :: rawBound :: raw := row | throw (IO.userError "expected word width, bound and source words")
  let bound := (← wordInput 64 rawBound).toNat
  let words ← raw.mapM fun w => BitVec.toNat <$> wordInput 64 w
  let width := width.toNat
  return String.intercalate "|" ([
    globalInt64N bound words, globalUint64N bound words, globalInt32N bound words,
    globalUint32N bound words, globalIntN width bound words, globalUintN width bound words,
    globalN true 8 bound words, globalN false 16 bound words].map wire)
