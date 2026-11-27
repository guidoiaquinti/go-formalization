import GoLean.Bytes.TrimFunc
import GoDiff.Math.Bits.Tuple

private def descriptor (kind : Nat) : Option (Nat → Bool) :=
  match kind with
  | 0 => none
  | 1 => some (fun _ => false)
  | 2 => some (fun _ => true)
  | 3 => some (· == 0xfffd)
  | 4 => some (fun r => r == 233 || r == 0x1f680)
  | _ => some (fun r => r % 7 == 2)

def main : IO Unit := runTuple fun row => do
  let kind :: spare :: raw := row | throw (IO.userError "expected predicate and input bytes")
  let bytes ← raw.mapM (wordInput 8)
  match GoLean.Bytes.trimFunc bytes spare.toNat (descriptor kind.toNat) with
  | none => return "panic"
  | some (none, calls) => return "nil|" ++ String.intercalate "," (calls.map toString)
  | some (some (remaining, capacity), calls) =>
    let output := String.intercalate "," (remaining.map (fun b => toString b.toNat))
    return toString capacity ++ ":" ++ output ++ "|" ++ String.intercalate "," (calls.map toString)
