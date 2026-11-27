import GoLean.Bytes.RepeatTitle
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def render (bytes : List Nat) (cap : Option Nat) : String :=
  let capacity := (cap.map toString).getD "growth"
  s!"slice:{capacity}:{String.intercalate "," (bytes.map toString)}"

def main : IO Unit := runTuple fun fields => do
  let maxInt :: count :: raw := fields | throw (IO.userError "expected count and bytes")
  let bytes ← raw.mapM (wordInput 8)
  let repeated := match repeatBytes (bytes.map BitVec.toNat) count maxInt.toNat with
    | .error .negative => "bytes: negative Repeat count"
    | .error .overflow => "bytes: Repeat output length overflow"
    | .ok output => render output (some output.length)
  let title := titleBytes bytes
  return repeated ++ "|" ++ render title.bytes title.exactCapacity
