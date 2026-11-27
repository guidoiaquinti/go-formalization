import GoLean.Strings.LastIndexFunc
import GoDiff.Math.Bits.Tuple

private def descriptor (kind : Nat) : Option (Nat → Bool) :=
  match kind with
  | 0 => none
  | 1 => some (fun _ => false)
  | 2 => some (fun _ => true)
  | 3 => some (· == 0xfffd)
  | 4 => some (· == 0x1f680)
  | _ => some (fun n => n % 7 == 2)

def main : IO Unit := runTuple fun row => do
  let kind :: raw := row | throw (IO.userError "expected predicate and bytes")
  let bytes ← raw.mapM (wordInput 8)
  match GoLean.Strings.lastIndexFunc bytes (descriptor kind.toNat) with
  | none => return "panic"
  | some (index, calls) =>
    return toString index ++ "|" ++ String.intercalate "," (calls.map toString)
