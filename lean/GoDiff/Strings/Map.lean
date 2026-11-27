import GoLean.Strings.Map
import GoDiff.Math.Bits.Tuple

private def descriptor (kind : Nat) (rune : Nat) : Int :=
  match kind with
  | 0 => rune
  | 1 => -1
  | 2 => if rune = 0xfffd then 120 else rune
  | 3 => if rune % 2 = 0 then 0x110000 else 0xd800
  | 4 => if rune % 2 = 0 then 0x1f680 else 0
  | _ => if rune % 3 = 0 then -2147483648 else 2147483647

def main : IO Unit := runTuple fun row => do
  let kind :: raw := row | throw (IO.userError "expected mapping and bytes")
  let bytes ← raw.mapM (wordInput 8)
  let (calls, out) := GoLean.Strings.mapRunes (descriptor kind.toNat) bytes
  return String.intercalate "," (calls.map toString) ++ "|" ++
    String.intercalate "," (out.map toString)
