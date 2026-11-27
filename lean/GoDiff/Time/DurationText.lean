import GoLean.Time.DurationText
import GoDiff.Math.Bits.Tuple

open GoLean.Time

private def wire (bytes : List (BitVec 8)) :=
  String.intercalate "," (bytes.map (toString ∘ BitVec.toNat))

def main : IO Unit := runTuple fun row => do
  let duration :: present :: raw := row | throw (IO.userError "expected duration and backing")
  let backing ← raw.mapM (wordInput 8)
  let rendered := wire ((durationString duration).toUTF8.data.toList.map UInt8.toBitVec)
  let privateResult := match durationFormat duration (if present == 1 then some backing else none) with
    | none => "nil"
    | some (start, storage) => s!"{start}:{wire storage}"
  return rendered ++ "|" ++ privateResult
