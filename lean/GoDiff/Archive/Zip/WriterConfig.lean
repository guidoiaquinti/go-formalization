import GoLean.Archive.Zip.WriterConfig
import GoDiff.Math.Bits.Tuple

open GoLean.Archive.Zip

private def wire (bytes : Comment) :=
  String.intercalate "," (bytes.map (toString ∘ BitVec.toNat))

private def report (result : Except ZipConfigFault ZipWriterConfig) (writer : Option ZipWriterConfig) :
    String :=
  let (status, observed) := match result with
    | .ok state => ("ok", some state)
    | .error .nilAccess => ("nil", writer)
    | .error .tooLong => ("tooLong", writer)
    | .error .afterWrite => ("afterWrite", writer)
  match observed with
  | none => status ++ ":none"
  | some w => s!"{status}:{w.count.map toString |>.getD "none"}:{wire w.comment}"

def main : IO Unit := runTuple fun row => do
  let tag :: count :: offset :: repeatCount :: fill :: oldLength :: raw := row |
    throw (IO.userError "expected writer state and comment")
  let fill ← wordInput 8 fill
  let raw ← raw.mapM (wordInput 8)
  let old := raw.take oldLength.toNat
  let comment := List.replicate repeatCount.toNat fill ++ raw.drop oldLength.toNat
  let writer := if tag == 0 then none else some ⟨old, if tag == 1 then none else some count⟩
  return report (writerSetComment writer comment) writer ++ "|" ++ report (writerSetOffset writer offset) writer
