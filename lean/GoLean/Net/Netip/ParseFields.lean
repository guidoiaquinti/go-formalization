namespace GoLean.Net.Netip

abbrev ParseBytes := List Nat

structure V4Fault where
  input : ParseBytes
  reason : String
  atBytes : ParseBytes

structure V4Fields where
  cells : ParseBytes
  fault : Option V4Fault

private def parseOctets (input : ParseBytes) : ParseBytes → ParseBytes → Nat → Nat → Nat → V4Fields
  | [], cells, pos, val, _ =>
    if pos < 3 then ⟨cells, some ⟨input, "IPv4 address too short", []⟩⟩
    else ⟨cells.set 3 val, none⟩
  | c :: tail, cells, pos, val, digits =>
    if 48 ≤ c && c ≤ 57 then
      let newVal := 10 * val + (c - 48)
      if digits == 1 && val == 0 then
        ⟨cells, some ⟨input, "IPv4 field has octet with leading zero", []⟩⟩
      else if newVal > 255 then
        ⟨cells, some ⟨input, "IPv4 field has value >255", []⟩⟩
      else parseOctets input tail cells pos newVal (digits + 1)
    else if c == 46 then
      if digits == 0 || tail.isEmpty then
        ⟨cells, some ⟨input, "IPv4 field must have at least one digit", c :: tail⟩⟩
      else if pos == 3 then
        ⟨cells, some ⟨input, "IPv4 address too long", []⟩⟩
      else parseOctets input tail (cells.set pos val) (pos + 1) 0 0
    else ⟨cells, some ⟨input, "unexpected character", c :: tail⟩⟩

def parseIPv4Fields (input : ParseBytes) (start stop : Nat) (fields : ParseBytes) : V4Fields :=
  parseOctets input ((input.drop start).take (stop - start)) fields 0 0 0

def parseIPv4 (input : ParseBytes) : Except V4Fault ParseBytes :=
  let result := parseIPv4Fields input 0 input.length [0, 0, 0, 0]
  match result.fault with
  | some fault => .error fault
  | none => .ok result.cells

structure AddrPortSplit where
  ip : ParseBytes
  port : ParseBytes
  bracketed : Bool

def splitAddrPort (input : ParseBytes) : Except String AddrPortSplit :=
  let reverse := input.reverse
  let port := reverse.takeWhile (· != 58)
  match reverse.dropWhile (· != 58) with
  | [] => .error "not an ip:port"
  | _ :: hostRev =>
    match hostRev.reverse with
    | [] => .error "no IP"
    | head :: rest =>
      if port.isEmpty then .error "no port"
      else if head == 91 then
        if rest.getLast? == some 93 then .ok ⟨rest.take (rest.length - 1), port.reverse, true⟩
        else .error "missing ]"
      else .ok ⟨head :: rest, port.reverse, false⟩

theorem parseIPv4Fields_partial :
    (parseIPv4Fields [49, 46, 50, 46, 120] 0 5 [42, 42, 42, 42]).cells = [1, 2, 42, 42] := rfl

theorem parseIPv4_allZero :
    parseIPv4 [48, 46, 48, 46, 48, 46, 48] = .ok [0, 0, 0, 0] := rfl

theorem splitAddrPort_emptyHostBrackets :
    splitAddrPort [91, 93, 58, 49] = .ok ⟨[], [49], true⟩ := rfl

end GoLean.Net.Netip
