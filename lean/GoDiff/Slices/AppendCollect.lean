import GoLean.Slices.AppendCollect
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

private def wire (cells : List Int) : String :=
  String.intercalate "," (cells.map toString)

private def describe (result : AppendSeqObservation Int Int) : String :=
  let status := match result.result with
    | .ok none => "nil"
    | .ok (some values) => "slice:" ++ wire values
    | .error .nilSequence => "nil-seq"
    | .error (.producer payload) => s!"producer:{payload}"
  s!"{status}/{result.acknowledged}/{result.reused}/{wire result.oldBacking}"

def main : IO Unit := runTuple fun fields => do
  let tag :: length :: capacity :: mode :: payload :: values := fields |
    throw (IO.userError "expected slice geometry, producer outcome and elements")
  if length < 0 || capacity < length || capacity.toNat > values.length then
    throw (IO.userError "invalid slice geometry")
  let initial := if tag == 0 then none else some (values.take length.toNat)
  let spare := (values.take capacity.toNat).drop length.toNat
  let sequence := if mode == 2 then none else
    some (values.drop capacity.toNat, if mode == 1 then Except.error payload else Except.ok ())
  return describe (sliceAppendSeq initial spare sequence) ++ "|" ++ describe (sliceCollect sequence)
