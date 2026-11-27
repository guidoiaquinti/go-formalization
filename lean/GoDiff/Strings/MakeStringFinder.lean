import GoLean.Strings.MakeStringFinder
import GoDiff.Math.Bits.Tuple

private def numbers (values : List Nat) : String :=
  String.intercalate "," (values.map toString)

def main : IO Unit := runTuple fun fields => do
  let count :: fields := fields | throw (IO.userError "expected pattern length and bytes")
  let pattern ← fields.mapM (wordInput 8)
  if count != Int.ofNat pattern.length then
    throw (IO.userError "pattern length differs from the supplied bytes")
  let finder := GoLean.Strings.makeStringFinder pattern
  return numbers (finder.pattern.map BitVec.toNat) ++ ":" ++
    numbers finder.badCharSkip ++ ":" ++ numbers finder.goodSuffixSkip
