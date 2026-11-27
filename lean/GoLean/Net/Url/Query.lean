import GoLean.Net.Url.Unescape
import GoLean.Net.Url.Values

namespace GoLean.Net.Url

private def signedDecimal (raw : URLBytes) (width : Nat) : Option Int := do
  let negative := raw.head? == some 45
  let signed := negative || raw.head? == some 43
  let digits := if signed then raw.drop 1 else raw
  if digits.isEmpty || !digits.all (fun b => 48 ≤ b.toNat && b.toNat ≤ 57) then none
  else
    let magnitude := digits.foldl (fun n b => 10 * n + (b.toNat - 48)) 0
    let half := 2 ^ (width - 1)
    if magnitude ≤ (if negative then half else half - 1) then
      some (if negative then -(Int.ofNat magnitude) else Int.ofNat magnitude)
    else none

/-- Go 1.27.0 src/net/url/url.go:963–977. The default comparison is
    inclusive. An accepted nonzero override uses strict comparison.
    The second result records the operational metric increment. -/
def urlParamsWithinMax (params : Int) (setting : URLBytes) (width : Nat) : Bool × Bool :=
  let standard := decide (params ≤ 10000)
  match signedDecimal setting width with
  | none => (standard, false)
  | some custom =>
    let within := custom == 0 || params < custom
    (within, standard != within)

inductive ParseQueryError where
  | limit
  | semicolon
  | unescape (error : UnescapeError)

private def cut (c : BitVec 8) (input : URLBytes) : URLBytes × URLBytes :=
  (input.takeWhile (· != c), (input.dropWhile (· != c)).drop 1)

private def parsePairs (query : URLBytes) (v : Values) (err : Option ParseQueryError) :
    Nat → Option (Values × Option ParseQueryError)
  | 0 => some (v, err)
  | fuel + 1 => do
    if query.isEmpty then return (v, err)
    let (field, next) := cut 38 query
    if field.contains 59 then parsePairs next v (some .semicolon) fuel
    else if field.isEmpty then parsePairs next v err fuel
    else
      let (key, value) := cut 61 field
      match do
        let k ← queryUnescape key
        let w ← queryUnescape value
        pure (k, w) with
      | .error failure =>
        parsePairs next v (if err.isSome then err else some (.unescape failure)) fuel
      | .ok (key, value) => do
        let nextMap ← valuesAdd v key value
        parsePairs next (some nextMap) err fuel

/-- Go 1.27.0 src/net/url/url.go:979–1011. Reject over-limit input
    before writes. Semicolons replace an earlier error, while invalid
    percent errors occupy only the first available error slot. -/
def parseQueryInto (v : Values) (query setting : URLBytes) (width : Nat) :
    Option (Values × Option ParseQueryError) :=
  let params := query.countP (· == 38) + 1
  if (urlParamsWithinMax (Int.ofNat params) setting width).1 then
    parsePairs query v none query.length
  else some (v, some .limit)

/-- Go 1.27.0 src/net/url/url.go:952–956. Start from a fresh nonnil map. -/
def parseQuery (query setting : URLBytes) (width : Nat) : Values × Option ParseQueryError :=
  (parseQueryInto (some []) query setting width).getD (some [], none)

/-- Go 1.27.0 src/net/url/url.go:1164–1167. Return partial values and discard the error. -/
def urlQuery (query : Option URLBytes) (setting : URLBytes) (width : Nat) : Option Values :=
  query.map (fun raw => (parseQuery raw setting width).1)

theorem urlParamsWithinMax_default :
    urlParamsWithinMax 10000 [] 64 = (true, false) ∧
    urlParamsWithinMax 10001 [] 64 = (false, false) := ⟨rfl, rfl⟩

theorem urlParamsWithinMax_exclusive :
    urlParamsWithinMax 5 [53] 64 = (false, true) ∧
    urlParamsWithinMax 10001 [48] 64 = (true, true) := ⟨rfl, rfl⟩

theorem parseQueryInto_limitBeforeWrite (v : Values) :
    parseQueryInto v [120] [49] 64 = some (v, some .limit) := rfl

theorem parseQueryInto_semicolonOverrides :
    parseQueryInto (some []) [37, 38, 59, 38, 120] [] 64 =
      some (some [([120], some [[]])], some .semicolon) := rfl

theorem parseQuery_empty : parseQuery [] [] 64 = (some [], none) := rfl

theorem parseQuery_order :
    parseQuery [97,61,49,38,97,61,50] [] 64 =
      (some [([97], some [[49],[50]])], none) := rfl

theorem urlQuery_nil (setting : URLBytes) (width : Nat) : urlQuery none setting width = none := rfl

theorem urlQuery_partial :
    urlQuery (some [97,38,37]) [] 64 = some (some [([97], some [[]])]) := rfl

end GoLean.Net.Url
