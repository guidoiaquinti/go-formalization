import GoLean.Net.Netip.ParseAddress
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Netip

private def wire (bytes : ParseBytes) : String := String.intercalate "," (bytes.map toString)

private def addrWire : ParsedAddr → String
  | .v4 bytes => "v4:" ++ wire bytes
  | .v6 bytes zone => "v6:" ++ wire bytes ++ ":" ++ wire zone

private def errWire (e : V4Fault) : String := s!"error:{e.reason}:{wire e.input}:{wire e.atBytes}"

private def parseWire : Except V4Fault ParsedAddr → String
  | .ok addr => addrWire addr
  | .error err => errWire err

def main : IO Unit := runTuple fun row => do
  let input := (row.drop 1).map Int.toNat
  let must := match mustParseAddr input with
    | .value addr => addrWire addr
    | .panic failure => "panic:" ++ errWire failure
  let text := match addrUnmarshalText input with
    | (some addr, _) => addrWire addr
    | (_, some err) => errWire err
    | (none, none) => "clear"
  return s!"{parseWire (parseIPv6 input)}|{parseWire (parseAddr input)}|{must}|{text}"
