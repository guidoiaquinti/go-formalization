import GoLean.Strings.NewReplacer
import GoDiff.Math.Bits.Tuple

private def readArguments : Nat → List Int → IO (List (List (BitVec 8)) × List Int)
  | 0, fields => pure ([], fields)
  | count + 1, length :: fields => do
    if length < 0 || length.toNat > fields.length then
      throw (IO.userError "invalid argument length")
    let arg ← (fields.take length.toNat).mapM (wordInput 8)
    let (args, remaining) ← readArguments count (fields.drop length.toNat)
    return (arg :: args, remaining)
  | _, _ => throw (IO.userError "missing argument length")

def main : IO Unit := runTuple fun row => do
  let count :: fields := row | throw (IO.userError "expected argument count")
  let (args, remaining) ← readArguments count.toNat fields
  if !remaining.isEmpty then throw (IO.userError "unexpected trailing bytes")
  match GoLean.Strings.newReplacer args with
  | none => return "strings.NewReplacer: odd argument count"
  | some state =>
    let fields := state.oldNew.getD [] |>.map fun arg =>
      s!"{arg.length}:" ++ String.intercalate "," (arg.map fun byte => toString byte.toNat)
    let storage := if state.oldNew.isNone then "nil" else "slice"
    return s!"{state.onceIsZero}:{state.implementationIsNil}:{storage}:" ++
      String.intercalate "|" fields
