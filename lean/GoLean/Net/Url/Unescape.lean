import GoLean.Net.Url.Escape

namespace GoLean.Net.Url

inductive UnescapeError where
  | escape (raw : List (BitVec 8))
  | host (byte : BitVec 8)

private def unescapedByte (high low : BitVec 8) : BitVec 8 :=
  (unhex high <<< 4) ||| unhex low

private def allowedTriple (mode : Encoding) (high low : BitVec 8) : Bool :=
  let pct := high == 50 && low == 53
  if mode == encodeHost then 8 ≤ (unhex high).toNat || pct
  else if mode == encodeZone then
    let decoded := unescapedByte high low
    pct || decoded == 32 || !shouldEscape decoded encodeHost
  else true

private def validate (mode : Encoding) : List (BitVec 8) →
    Except UnescapeError (Nat × Bool)
  | [] => .ok (0, false)
  | byte :: rest =>
    if byte == 37 then
      match rest with
      | high :: low :: remaining => do
        if isHex high && isHex low && allowedTriple mode high low then
          let (count, plus) ← validate mode remaining
          return (count + 1, plus)
        else .error (.escape [byte, high, low])
      | _ => .error (.escape (byte :: rest))
    else if byte == 43 then do
      let (count, plus) ← validate mode rest
      return (count, plus || mode == encodeQueryComponent)
    else if (mode == encodeHost || mode == encodeZone) && byte.toNat < 128 && shouldEscape byte mode then
      .error (.host byte)
    else validate mode rest

private def decode (mode : Encoding) : List (BitVec 8) → List (BitVec 8)
  | [] => []
  | byte :: rest =>
    if byte == 37 then
      match rest with
      | high :: low :: remaining => unescapedByte high low :: decode mode remaining
      | _ => []
    else if byte == 43 && mode == encodeQueryComponent then
      32 :: decode mode rest
    else byte :: decode mode rest

/-- Go 1.27.0 src/net/url/url.go:106–182. Validate from left to right
    before decoding, preserving exact payload bytes on the first error.
    The two Go error types are represented by separate constructors. -/
def unescape (input : List (BitVec 8)) (mode : Encoding) :
    Except UnescapeError (List (BitVec 8)) := do
  let (count, hasPlus) ← validate mode input
  if count == 0 && !hasPlus then return input
  else return decode mode input

def queryUnescape (input : List (BitVec 8)) : Except UnescapeError (List (BitVec 8)) :=
  unescape input encodeQueryComponent

def pathUnescape (input : List (BitVec 8)) : Except UnescapeError (List (BitVec 8)) :=
  unescape input encodePathSegment

theorem queryUnescape_plus : queryUnescape [43, 37, 50, 66] = .ok [32, 43] := rfl

theorem queryUnescape_firstInvalid (rest : List (BitVec 8)) :
    queryUnescape ([43, 37, 71, 49] ++ rest) = .error (.escape [37, 71, 49]) := rfl

theorem pathUnescape_plus : pathUnescape [43, 37, 50, 66] = .ok [43, 43] := rfl

theorem pathUnescape_truncated : pathUnescape [97, 37, 70] = .error (.escape [37, 70]) := rfl

theorem unescape_host_percent :
    unescape [37, 50, 53] encodeHost = .ok [37] ∧
    unescape [37, 50, 48] encodeHost = .error (.escape [37, 50, 48]) := ⟨rfl, rfl⟩

theorem unescape_zone_space :
    unescape [37, 50, 48] encodeZone = .ok [32] ∧
    unescape [32] encodeZone = .error (.host 32) := ⟨rfl, rfl⟩

theorem unescape_host_errorPrefix (rest : List (BitVec 8)) :
    unescape (47 :: rest) encodeHost = .error (.host 47) := by
  cases rest with
  | nil => rfl
  | cons first tail => cases tail <;> rfl

end GoLean.Net.Url
