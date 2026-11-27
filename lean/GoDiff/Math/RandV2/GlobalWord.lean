import GoLean.Math.RandV2.GlobalWord
import GoDiff.Math.Bits.Tuple

open GoLean.Math.RandV2

private def wire (r : WordResult) := s!"{r.status}:{r.value}:{r.pulls}"

def main : IO Unit := runTuple fun row => do
  let [width, mode, value] := row | throw (IO.userError "expected word width, mode and value")
  let word ← wordInput 64 value
  let source := if mode == 2 then none else if mode == 1 then some WordSource.raises
    else some (.word word.toNat)
  let r := newRand source
  let width := width.toNat
  let all := [globalInt64 r, globalUint32 r, globalUint64 r, globalInt32 r,
    globalInt width r, globalUint width r, globalFloat64 r, globalFloat32 r]
  return String.intercalate "|" (all.map wire)
