import GoLean.Net.Netip.ParseFields
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Netip

private def wire (xs : ParseBytes) : String := String.intercalate "," (xs.map toString)

private def faultWire (err : V4Fault) : String := s!"{err.reason}:{wire err.input}:{wire err.atBytes}"

def main : IO Unit := runTuple fun fields => do
  let start :: stop :: a :: b :: c :: d :: bytes := fields |
    throw (IO.userError "expected slice bounds, four initial fields, and input bytes")
  let input := bytes.map Int.toNat
  let part := parseIPv4Fields input start.toNat stop.toNat ([a,b,c,d].map Int.toNat)
  let partStatus := (part.fault.map faultWire).getD "ok"
  let whole := match parseIPv4 input with
    | .ok cells => s!"ok:{wire cells}"
    | .error err => faultWire err
  let port := match splitAddrPort input with
    | .ok result => s!"{wire result.ip}:{wire result.port}:{result.bracketed}"
    | .error message => message
  return s!"{wire part.cells}|{partStatus}|{whole}|{port}"
