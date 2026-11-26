import GoLean.Strconv.ParseBool
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let 0 :: input := fields | throw (IO.userError "expected marker followed by input bytes")
  unless input.all (fun x => 0 ≤ x && x < 256) do
    throw (IO.userError "invalid input byte")
  let (value, error) := GoLean.Strconv.parseBool (input.map fun x => x.toNat.toUInt8)
  match error with
  | none => return s!"{value} ok"
  | some e => return s!"{value} syntax {e.func} [{String.intercalate " " (e.num.map toString)}]"
