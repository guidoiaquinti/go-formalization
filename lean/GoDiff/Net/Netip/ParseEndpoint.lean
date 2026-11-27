import GoLean.Net.Netip.ParseEndpoint
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Netip

private def wire (cells : ParseBytes) := String.intercalate "," (cells.map toString)

private def valueWire (addr : ParsedAddr) (number : Nat) : String :=
  match addr with
  | .v4 bytes => s!"v4:{wire bytes}::{number}"
  | .v6 bytes zone => s!"v6:{wire bytes}:{wire zone}:{number}"

private def parseWire : Except String (ParsedAddr × Nat) → String
  | .ok (a, n) => valueWire a n
  | .error reason => "error:" ++ reason

private def mustWire : MustEndpoint → String
  | .value a n => valueWire a n
  | .panic reason => "panic:" ++ reason

private def textWire : Option (ParsedAddr × Nat) × Option String → String
  | (some (a, n), _) => valueWire a n
  | (_, some reason) => "error:" ++ reason
  | (none, none) => "clear"

def main : IO Unit := runTuple fun row => do
  let width :: raw := row | throw (IO.userError "expected native size and input bytes")
  let input := raw.map Int.toNat
  return String.intercalate "|" [parseWire (parseAddrPort input), mustWire (mustParseAddrPort input),
    textWire (addrPortUnmarshalText input), parseWire (parsePrefix input width.toNat),
    mustWire (mustParsePrefix input width.toNat), textWire (prefixUnmarshalText input width.toNat)]
