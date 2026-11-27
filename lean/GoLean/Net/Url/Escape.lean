import Init.Data.BitVec.Basic

namespace GoLean.Net.Url

abbrev Encoding := BitVec 8

def encodePath : Encoding := 1
def encodePathSegment : Encoding := 2
def encodeHost : Encoding := 4
def encodeZone : Encoding := 8
def encodeUserPassword : Encoding := 16
def encodeQueryComponent : Encoding := 32
def encodeFragment : Encoding := 64

private def between (byte lo hi : Nat) : Bool := lo ≤ byte && byte ≤ hi

private def alphaNum (byte : Nat) : Bool :=
  between byte 48 57 || between byte 65 90 || between byte 97 122

private def hasHexBit (byte : Nat) : Bool :=
  between byte 48 57 || between byte 65 70 || between byte 97 102

/-- Reconstruct Go 1.27.0 encoding_table.go from the generator's allowed
    byte sets. The high bit serves ishex and participates in arbitrary modes. -/
def encodingByte (c : BitVec 8) : Encoding :=
  let byte := c.toNat
  let common := alphaNum byte || [45, 95, 46, 126].contains byte
  let allowed (chars : List Nat) := common || chars.contains byte
  let host := allowed [33, 36, 38, 39, 40, 41, 42, 43, 44, 59, 61, 58, 91, 93, 60, 62, 34]
  (if allowed [36, 38, 43, 44, 47, 58, 59, 61, 64] then encodePath else 0) |||
  (if allowed [36, 38, 43, 58, 61, 64] then encodePathSegment else 0) |||
  (if host then encodeHost ||| encodeZone else 0) |||
  (if allowed [36, 38, 43, 44, 59, 61] then encodeUserPassword else 0) |||
  (if common then encodeQueryComponent else 0) |||
  (if allowed [36, 38, 43, 44, 47, 58, 59, 61, 63, 64, 33, 40, 41, 42] then encodeFragment else 0) |||
  (if hasHexBit byte then 128 else 0)

def shouldEscape (c : BitVec 8) (mode : Encoding) : Bool :=
  (encodingByte c &&& mode) == 0

def isHex (c : BitVec 8) : Bool :=
  (encodingByte c &&& 128) != 0

def unhex (c : BitVec 8) : BitVec 8 :=
  9 * (c >>> 6) + (c &&& 15)

private def upperHex (n : Nat) : BitVec 8 :=
  BitVec.ofNat 8 (if n < 10 then n + 48 else n + 55)

private def isQuerySpace (mode : Encoding) (c : BitVec 8) : Bool :=
  mode == encodeQueryComponent && c == 32

private def writeEscaped (mode : Encoding) (c : BitVec 8) : List (BitVec 8) :=
  if isQuerySpace mode c then [43]
  else if shouldEscape c mode then [37, upperHex (c.toNat / 16), upperHex (c.toNat % 16)]
  else [c]

/-- Go 1.27.0 src/net/url/url.go:196–247. Count escaped bytes before
    copying and use uppercase nibble encoding. Storage location is abstracted. -/
def escape (input : List (BitVec 8)) (mode : Encoding) : List (BitVec 8) :=
  let escaped := input.filter (fun c => shouldEscape c mode)
  let spaces := escaped.countP (isQuerySpace mode)
  let hexes := escaped.length - spaces
  if spaces = 0 && hexes = 0 then input
  else if hexes = 0 then input.map (fun c => if c == 32 then 43 else c)
  else input.flatMap (writeEscaped mode)

def queryEscape (input : List (BitVec 8)) : List (BitVec 8) :=
  escape input encodeQueryComponent

def pathEscape (input : List (BitVec 8)) : List (BitVec 8) :=
  escape input encodePathSegment

theorem shouldEscape_zero (c : BitVec 8) : shouldEscape c 0 = true := by
  simp [shouldEscape]

theorem shouldEscape_combined :
    shouldEscape 47 encodePath = false ∧ shouldEscape 47 encodePathSegment = true ∧
    shouldEscape 47 (encodePath ||| encodePathSegment) = false := ⟨rfl, rfl, rfl⟩

theorem isHex_examples : isHex 48 = true ∧ isHex 102 = true ∧ isHex 103 = false := ⟨rfl, rfl, rfl⟩

theorem unhex_examples : unhex 57 = 9 ∧ unhex 65 = 10 ∧ unhex 102 = 15 := ⟨rfl, rfl, rfl⟩

theorem escape_empty (mode : Encoding) : escape [] mode = [] := rfl

theorem escape_high : escape [255, 0, 47] encodeFragment = [37, 70, 70, 37, 48, 48, 47] := rfl

theorem queryEscape_plus : queryEscape [32, 43] = [43, 37, 50, 66] := rfl

theorem pathEscape_plus : pathEscape [32, 43, 47] = [37, 50, 48, 43, 37, 50, 70] := rfl

end GoLean.Net.Url
