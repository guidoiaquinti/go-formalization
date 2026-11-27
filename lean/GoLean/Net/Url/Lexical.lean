import GoLean.Net.Url.Values

namespace GoLean.Net.Url

private def decimalByte (c : BitVec 8) : Bool := 48 ≤ c.toNat && c.toNat ≤ 57

private def alphaByte (c : BitVec 8) : Bool :=
  (65 ≤ c.toNat && c.toNat ≤ 90) || (97 ≤ c.toNat && c.toNat ≤ 122)

/-- Go 1.27.0 src/net/url/url.go:762–775. A colon alone is valid.
    Digits have no numerical conversion or range restriction. -/
def validOptionalPort : URLBytes → Bool
  | [] => true
  | c :: rest => c == 58 && rest.all decimalByte

private def stripBrackets (host : URLBytes) : URLBytes :=
  if host.head? == some 91 && host.getLast? == some 93 then
    (host.drop 1).take (host.length - 2)
  else host

/-- Go 1.27.0 src/net/url/url.go:1209–1223. Split on the last colon
    only if its suffix is a valid optional port. Strip outer brackets
    from the resulting host independently of port extraction. -/
def splitHostPort (input : URLBytes) : URLBytes × URLBytes :=
  let candidate := input.reverse.findIdx? (· == 58)
  let (host, port) := match candidate with
    | none => (input, [])
    | some fromEnd =>
      let colon := input.length - fromEnd - 1
      if validOptionalPort (input.drop colon) then (input.take colon, input.drop (colon + 1))
      else (input, [])
  (stripBrackets host, port)

/-- Go 1.27.0 src/net/url/url.go:1193–1197. -/
def urlHostname (host : Option URLBytes) : Option URLBytes :=
  host.map (fun input => (splitHostPort input).1)

/-- Go 1.27.0 src/net/url/url.go:1201–1205. -/
def urlPort (host : Option URLBytes) : Option URLBytes :=
  host.map (fun input => (splitHostPort input).2)

/-- Go 1.27.0 src/net/url/url.go:1101–1103. These accessors represent a
    nil receiver fault by an absent result. -/
def urlIsAbs (scheme : Option URLBytes) : Option Bool :=
  scheme.map (fun input => !input.isEmpty)

private def collectScheme (original : URLBytes) :
    URLBytes → URLBytes → Option (URLBytes × URLBytes)
  | _, [] => some ([], original)
  | reversed, c :: rest =>
    if c == 58 then
      if reversed.isEmpty then none else some (reversed.reverse, rest)
    else if alphaByte c || (!reversed.isEmpty && (decimalByte c || [43, 45, 46].contains c)) then
      collectScheme original (c :: reversed) rest
    else some ([], original)

/-- Go 1.27.0 src/net/url/url.go:369–393. Returns an unchanged scheme
    spelling and remaining suffix, or a missing-protocol-scheme error. -/
def getScheme (raw : URLBytes) : Option (URLBytes × URLBytes) :=
  collectScheme raw [] raw

/-- Go 1.27.0 src/net/url/url.go:1317–1325. Detect raw C0 and DEL. -/
def stringContainsCTLByte (raw : URLBytes) : Bool :=
  raw.any fun c => c.toNat < 32 || c == 127

theorem validOptionalPort_colon : validOptionalPort [58] = true := rfl
theorem validOptionalPort_plus : validOptionalPort [58, 43, 49] = false := rfl

theorem splitHostPort_brackets :
    splitHostPort [91, 58, 58, 49, 93, 58, 48, 56, 48] = ([58, 58, 49], [48, 56, 48]) := rfl

theorem splitHostPort_badPort :
    splitHostPort [91, 58, 58, 49, 93, 58, 120] = ([91, 58, 58, 49, 93, 58, 120], []) := rfl

theorem urlHostname_nil : urlHostname none = none := rfl
theorem urlHostname_unported : urlHostname (some [91, 58, 58, 49, 93]) = some [58, 58, 49] := rfl
theorem urlPort_nil : urlPort none = none := rfl
theorem urlPort_bad : urlPort (some [97, 58, 120]) = some [] := rfl
theorem urlIsAbs_nil : urlIsAbs none = none := rfl
theorem urlIsAbs_nonempty : urlIsAbs (some [0]) = some true := rfl

theorem getScheme_missing (rest : URLBytes) : getScheme (58 :: rest) = none := rfl
theorem getScheme_preservesCase :
    getScheme [65, 43, 49, 58, 255] = some ([65, 43, 49], [255]) := rfl

theorem stringContainsCTLByte_del : stringContainsCTLByte [127] = true := rfl
theorem stringContainsCTLByte_high : stringContainsCTLByte [128, 255] = false := rfl

end GoLean.Net.Url
