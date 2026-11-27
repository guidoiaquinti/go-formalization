import GoLean.Net.Url.Unescape

namespace GoLean.Net.Url

structure URLPart where
  decoded : List (BitVec 8)
  raw : List (BitVec 8)

/-- Go 1.27.0 src/net/url/url.go:703–724. This byte filter accepts
    a percent byte without validating the percent escape it starts. -/
def validEncoded (input : List (BitVec 8)) (mode : Encoding) : Bool :=
  input.all fun byte =>
    [33, 36, 38, 39, 40, 41, 42, 43, 44, 59, 61, 58, 64, 91, 93, 37].contains byte ||
      !shouldEscape byte mode

private def keepRaw (part : URLPart) (mode : Encoding) : Bool :=
  !part.raw.isEmpty && validEncoded part.raw mode &&
    match unescape part.raw mode with
    | .ok candidate => candidate == part.decoded
    | .error _ => false

private def escapedPart (receiver : Option URLPart) (mode : Encoding) (star : Bool) :
    Option (List (BitVec 8)) :=
  receiver.map fun part =>
    if keepRaw part mode then part.raw
    else if star && part.decoded == [42] then [42]
    else escape part.decoded mode

/-- Go 1.27.0 src/net/url/url.go:687–699. Raw-path validation precedes
    the asterisk special case. None models a nil receiver panic. -/
def escapedPath (receiver : Option URLPart) : Option (List (BitVec 8)) :=
  escapedPart receiver encodePath true

/-- Go 1.27.0 src/net/url/url.go:750–759. No asterisk special case. -/
def escapedFragment (receiver : Option URLPart) : Option (List (BitVec 8)) :=
  escapedPart receiver encodeFragment false

/-- None is a nil-receiver panic. A validation error retains the old state,
    even when the receiver was nil and therefore cannot be updated. -/
private def setPart (receiver : Option URLPart) (raw : List (BitVec 8)) (mode : Encoding) :
    Option (Except UnescapeError Unit × Option URLPart) :=
  match unescape raw mode with
  | .error err => some (.error err, receiver)
  | .ok decoded => receiver.map fun _ =>
    (.ok (), some ⟨decoded, if raw == escape decoded mode then [] else raw⟩)

/-- Go 1.27.0 src/net/url/url.go:660–672. -/
def setPath (receiver : Option URLPart) (raw : List (BitVec 8)) :
    Option (Except UnescapeError Unit × Option URLPart) :=
  setPart receiver raw encodePath

/-- Go 1.27.0 src/net/url/url.go:727–739. -/
def setFragment (receiver : Option URLPart) (raw : List (BitVec 8)) :
    Option (Except UnescapeError Unit × Option URLPart) :=
  setPart receiver raw encodeFragment

theorem validEncoded_percent (mode : Encoding) : validEncoded [37] mode = true := rfl

theorem validEncoded_space : validEncoded [32] encodePath = false := rfl

theorem escapedPath_rawStar :
    escapedPath (some ⟨[42], [37, 50, 97]⟩) = some [37, 50, 97] := rfl

theorem escapedPath_defaultStar : escapedPath (some ⟨[42], []⟩) = some [42] := rfl

theorem escapedFragment_mismatch :
    escapedFragment (some ⟨[32], [37, 50, 49]⟩) = some [37, 50, 48] := rfl

theorem setPath_errorPreserves (receiver : Option URLPart) :
    setPath receiver [37, 49] = some (.error (.escape [37, 49]), receiver) := rfl

theorem setPath_nonCanonical (receiver : URLPart) :
    setPath (some receiver) [37, 50, 102] = some (.ok (), some ⟨[47], [37, 50, 102]⟩) := rfl

theorem setFragment_default (receiver : URLPart) :
    setFragment (some receiver) [37, 50, 48] = some (.ok (), some ⟨[32], []⟩) := rfl

theorem setFragment_nilInvalid :
    setFragment none [37] = some (.error (.escape [37]), none) := rfl

end GoLean.Net.Url
