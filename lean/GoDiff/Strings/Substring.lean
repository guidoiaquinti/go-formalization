import GoLean.Strings.LastIndex
import GoDiff.Math.Bits.Tuple

private def encode (words : List Nat) : String :=
  String.intercalate " " (toString words.length :: words.map toString)

def main : IO Unit := runTuple fun fields => do
  let n :: words := fields | throw (IO.userError "expected needle length and bytes")
  if n < 0 || words.length < n.toNat then
    throw (IO.userError "invalid needle length")
  let bytes ← words.mapM fun word => do
    return (← wordInput 8 word).toNat
  let needle := bytes.take n.toNat
  let input := bytes.drop n.toNat
  let found := (GoLean.Strings.index input needle).map toString |>.getD "-1"
  let (before, after, success) := GoLean.Strings.cut input needle
  let last := (GoLean.Strings.lastIndex input needle).map toString |>.getD "-1"
  let (lastBefore, lastAfter, lastFound) := GoLean.Strings.cutLast input needle
  return s!"{found} {last} {GoLean.Strings.contains input needle} {success} {encode before} {encode after} {lastFound} {encode lastBefore} {encode lastAfter}"
