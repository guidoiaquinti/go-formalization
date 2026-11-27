import GoLean.Net.Netip.ParseFields

namespace GoLean.Net.Netip

inductive ParsedAddr where
  | v4 (bytes : ParseBytes)
  | v6 (bytes zone : ParseBytes)

private def hex (c : Nat) : Option Nat :=
  if 48 ≤ c && c ≤ 57 then some (c - 48)
  else if 65 ≤ c && c ≤ 70 then some (c - 55)
  else if 97 ≤ c && c ≤ 102 then some (c - 87)
  else none

private def scanGroup : ParseBytes → Nat → Nat → Nat × Nat × ParseBytes
  | [], count, acc => (count, acc, [])
  | b :: tail, count, acc =>
    match hex b with
    | none => (count, acc, b :: tail)
    | some digit =>
      if count == 4 then (5, acc, b :: tail)
      else scanGroup tail (count + 1) (16 * acc + digit)

private def v6Error (input : ParseBytes) (why : String) (suffix : ParseBytes := []) :
    Except V4Fault ParsedAddr := .error ⟨input, why, suffix⟩

private def finishV6 (input bytes zone : ParseBytes) (ellipse : Option Nat) :
    Except V4Fault ParsedAddr :=
  if bytes.length < 16 then
    match ellipse with
    | none => v6Error input "address string too short"
    | some offset =>
      .ok (.v6 (bytes.take offset ++ List.replicate (16 - bytes.length) 0 ++ bytes.drop offset) zone)
  else if ellipse.isSome then v6Error input "the :: must expand to at least one field of zeros"
  else .ok (.v6 bytes zone)

private def ipv6Groups (input zone : ParseBytes) (endAt : Nat) :
    Nat → ParseBytes → ParseBytes → Option Nat → Except V4Fault ParsedAddr
  | 0, rest, bytes, ellipse =>
    if rest.isEmpty then finishV6 input bytes zone ellipse
    else v6Error input "trailing garbage after address" rest
  | fuel + 1, rest, bytes, ellipse =>
    if bytes.length == 16 then
      if rest.isEmpty then finishV6 input bytes zone ellipse
      else v6Error input "trailing garbage after address" rest
    else
      let (n, acc, suffix) := scanGroup rest 0 0
      if n > 4 then v6Error input "each group must have 4 or less digits" rest
      else if n == 0 then
        v6Error input "each colon-separated field must have at least one digit" rest
      else if suffix.head? == some 46 then
        if ellipse.isNone && bytes.length != 12 then
          v6Error input "embedded IPv4 address must replace the final 2 fields of the address" rest
        else if bytes.length + 4 > 16 then
          v6Error input "too many hex fields to fit an embedded IPv4 at the end of the address" rest
        else
          let v4 := parseIPv4Fields input (endAt - rest.length) endAt [0, 0, 0, 0]
          match v4.fault with
          | some failure => .error failure
          | none => finishV6 input (bytes ++ v4.cells) zone ellipse
      else
        let bytes := bytes ++ [acc / 256, acc % 256]
        match suffix with
        | [] => finishV6 input bytes zone ellipse
        | 58 :: [] => v6Error input "colon must be followed by more characters" suffix
        | 58 :: 58 :: rest =>
          if ellipse.isSome then v6Error input "multiple :: in address" (58 :: rest)
          else if rest.isEmpty then finishV6 input bytes zone (some bytes.length)
          else ipv6Groups input zone endAt fuel rest bytes (some bytes.length)
        | 58 :: rest => ipv6Groups input zone endAt fuel rest bytes ellipse
        | _ => v6Error input "unexpected character, want colon" suffix

def parseIPv6 (input : ParseBytes) : Except V4Fault ParsedAddr :=
  let body := input.takeWhile (· != 37)
  let delim := input.drop body.length
  let zone := delim.drop 1
  if !delim.isEmpty && zone.isEmpty then v6Error input "zone must be a non-empty string"
  else match body with
    | [58, 58] => .ok (.v6 (List.replicate 16 0) zone)
    | 58 :: 58 :: rest => ipv6Groups input zone body.length 8 rest [] (some 0)
    | _ => ipv6Groups input zone body.length 8 body [] none

def parseAddr (input : ParseBytes) : Except V4Fault ParsedAddr :=
  let separator := input.find? fun b => b == 46 || b == 58 || b == 37
  match separator with
  | some 46 => (parseIPv4 input).map ParsedAddr.v4
  | some 58 => parseIPv6 input
  | some 37 => v6Error input "missing IPv6 address"
  | _ => v6Error input "unable to parse IP"

inductive MustAddr where
  | value (addr : ParsedAddr)
  | panic (fault : V4Fault)

def mustParseAddr (input : ParseBytes) : MustAddr :=
  match parseAddr input with
  | .ok addr => .value addr
  | .error err => .panic err

def addrUnmarshalText (input : ParseBytes) : Option ParsedAddr × Option V4Fault :=
  if input.isEmpty then (none, none)
  else match parseAddr input with
    | .ok addr => (some addr, none)
    | .error err => (none, some err)

theorem parseIPv6_unspecified : parseIPv6 [58, 58] = .ok (.v6 (List.replicate 16 0) []) := rfl

theorem parseAddr_zoneFirst : parseAddr [37, 58, 58] =
    .error ⟨[37, 58, 58], "missing IPv6 address", []⟩ := rfl

theorem mustParseAddr_empty : mustParseAddr [] = .panic ⟨[], "unable to parse IP", []⟩ := rfl

theorem addrUnmarshalText_empty : addrUnmarshalText [] = (none, none) := rfl

end GoLean.Net.Netip
