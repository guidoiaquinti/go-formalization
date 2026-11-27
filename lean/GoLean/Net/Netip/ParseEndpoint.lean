import GoLean.Net.Netip.ParseAddress

namespace GoLean.Net.Netip

private def decimal (bytes : ParseBytes) (max : Nat) : Option Nat :=
  if bytes.isEmpty then none else
    bytes.foldlM (m := Option) (fun acc c =>
      if c < 48 || c > 57 then none else
        let next := acc * 10 + (c - 48)
        if next > max then none else some next) 0

private def isV6 : ParsedAddr → Bool
  | .v4 _ => false
  | .v6 _ _ => true

def parseAddrPort (input : ParseBytes) : Except String (ParsedAddr × Nat) := do
  let parts ← splitAddrPort input
  let some port := decimal parts.port 65535 | throw "invalid port"
  let addr ← (parseAddr parts.ip).mapError (fun _ => "address fault")
  if parts.bracketed && !isV6 addr then throw "IPv4 bracketed"
  else if !parts.bracketed && isV6 addr then throw "IPv6 unbracketed"
  else return (addr, port)

def parsePrefix (input : ParseBytes) (wordSize : Nat) : Except String (ParsedAddr × Nat) := do
  let suffix := input.reverse.takeWhile (· != 47)
  let _ :: host := input.reverse.drop suffix.length | throw "no slash"
  let addr ← (parseAddr host.reverse).mapError (fun _ => "address fault")
  if let .v6 _ (_ :: _) := addr then throw "zone"
  let bits := suffix.reverse
  if bits.length > 1 && (bits.head! < 49 || bits.head! > 57) then throw "bad bits"
  let some len := decimal bits (2^(wordSize-1)-1) | throw "bad bits"
  if len > (if isV6 addr then 128 else 32) then throw "out of range"
  return (addr, len)

inductive MustEndpoint where
  | value (addr : ParsedAddr) (number : Nat)
  | panic (reason : String)

private def mustEndpoint : Except String (ParsedAddr × Nat) → MustEndpoint
  | .ok (addr, number) => .value addr number
  | .error reason => .panic reason

private def unmarshalEndpoint (input : ParseBytes) (parsed : Except String (ParsedAddr × Nat)) :
    Option (ParsedAddr × Nat) × Option String :=
  if input.isEmpty then (none, none)
  else match parsed with
    | .ok value => (some value, none)
    | .error reason => (none, some reason)

def mustParseAddrPort (input : ParseBytes) : MustEndpoint := mustEndpoint (parseAddrPort input)

def mustParsePrefix (input : ParseBytes) (wordSize : Nat) : MustEndpoint :=
  mustEndpoint (parsePrefix input wordSize)

def addrPortUnmarshalText (input : ParseBytes) : Option (ParsedAddr × Nat) × Option String :=
  unmarshalEndpoint input (parseAddrPort input)

def prefixUnmarshalText (input : ParseBytes) (wordSize : Nat) :
    Option (ParsedAddr × Nat) × Option String :=
  unmarshalEndpoint input (parsePrefix input wordSize)

theorem parseAddrPort_zero : parseAddrPort [91,58,58,93,58,48] =
    .ok (.v6 (List.replicate 16 0) [], 0) := rfl

theorem parsePrefix_noSlash (width : Nat) : parsePrefix [] width = .error "no slash" := rfl

theorem mustParseAddrPort_empty : mustParseAddrPort [] = .panic "not an ip:port" := rfl

theorem mustParsePrefix_empty (width : Nat) : mustParsePrefix [] width = .panic "no slash" := rfl

theorem addrPortUnmarshalText_empty : addrPortUnmarshalText [] = (none, none) := rfl

theorem prefixUnmarshalText_empty (width : Nat) : prefixUnmarshalText [] width = (none, none) := rfl

end GoLean.Net.Netip
