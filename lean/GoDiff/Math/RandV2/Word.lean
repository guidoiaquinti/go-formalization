import GoLean.Math.RandV2.Word
import GoDiff.Math.Bits.Tuple

open GoLean.Math.RandV2

private def wire (r : WordResult) := s!"{r.status}:{r.value}:{r.pulls}"

def main : IO Unit := runTuple fun row => do
  let [width, mode, value] := row | throw (IO.userError "expected word width, mode and value")
  let word ← wordInput 64 value
  let source := if mode == 2 then none else if mode == 1 then some WordSource.raises
    else some (.word word.toNat)
  let r := newRand source
  let receiver := if mode == 3 then none else some r
  let width := width.toNat
  let all := [randInt64 receiver, randUint32 receiver, randUint64 receiver, randInt32 receiver,
    randInt width receiver, randUint width receiver, randFloat64 receiver, randFloat32 receiver]
  let installed := if r.source.isNone then "nil" else "source"
  return installed ++ ":0|" ++ String.intercalate "|" (all.map wire)
