import GoLean.Strings.ToTitleSpecial
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

private def caseRows : Nat → List Int → IO (List SpecialCaseRange × List Int)
  | 0, fields => pure ([], fields)
  | count + 1, lo :: hi :: upper :: lower :: title :: fields => do
    let entry : SpecialCaseRange :=
      ⟨BitVec.ofInt 32 lo, BitVec.ofInt 32 hi, BitVec.ofInt 32 upper,
        BitVec.ofInt 32 lower, BitVec.ofInt 32 title⟩
    let (entries, raw) ← caseRows count fields
    return (entry :: entries, raw)
  | _, _ => throw (IO.userError "expected range bounds and case deltas")

def main : IO Unit := runTuple fun fields => do
  let count :: fields := fields | throw (IO.userError "expected case table size")
  let (entries, raw) ← caseRows count.toNat fields
  let input ← raw.mapM (wordInput 8)
  let result := toTitleSpecial entries.toArray input
  return "bytes:" ++ String.intercalate "," (result.map toString)
