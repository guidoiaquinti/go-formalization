import GoLean.Net.Url.PathFragment
import GoLean.Net.Url.Userinfo

namespace GoLean.Net.Url

structure URLRender where
  scheme : URLBytes := []
  opaquePart : URLBytes := []
  credentials : Option Userinfo := none
  host : URLBytes := []
  path : URLPart := ⟨[], []⟩
  omitHost : Bool := false
  forceQuery : Bool := false
  query : URLBytes := []
  fragment : URLPart := ⟨[], []⟩

private def literal (s : String) : URLBytes :=
  s.toUTF8.data.toList.map UInt8.toBitVec

private def querySuffix (u : URLRender) : URLBytes :=
  if u.forceQuery || !u.query.isEmpty then 63 :: u.query else []

private def startsSlash (path : URLBytes) : Bool := path.head? == some 47

private def startsDoubleSlash (path : URLBytes) : Bool :=
  path.take 2 == [47, 47]

private def renderPath (u : URLRender) (scheme : URLBytes) : URLBytes :=
  let bare := u.host.isEmpty && u.credentials.isNone
  let omitAuthority := u.omitHost && bare
  let head :=
    if omitAuthority || (u.scheme.isEmpty && bare) then scheme
    else
      let slash := if bare && u.path.decoded.isEmpty then [] else [47, 47]
      let credentials := if u.credentials.isSome then userinfoString u.credentials ++ [64] else []
      scheme ++ slash ++ credentials ++ escape u.host encodeHost
  let path := (escapedPath (some u.path)).getD []
  if omitAuthority && startsDoubleSlash path then head ++ literal "%2F" ++ path.drop 1
  else
    let slash := if !path.isEmpty && !startsSlash path && !u.host.isEmpty then [47] else []
    let dot := if head.isEmpty && (path.takeWhile (· != 47)).contains 58 then [46, 47] else []
    head ++ slash ++ dot ++ path

/-- Go 1.27.0 url.go:798–873. Nil signals a receiver access failure. -/
def urlString (receiver : Option URLRender) : Option URLBytes :=
  receiver.map fun u =>
    let scheme := if u.scheme.isEmpty then [] else u.scheme ++ [58]
    let main := if u.opaquePart.isEmpty then renderPath u scheme else scheme ++ u.opaquePart
    let fragment := if u.fragment.decoded.isEmpty then []
      else 35 :: (escapedFragment (some u.fragment)).getD []
    main ++ querySuffix u ++ fragment

def urlRedacted (receiver : Option URLRender) : URLBytes :=
  let changed := receiver.map fun u =>
    { u with credentials := u.credentials.map (fun info =>
      if info.passwordSet then { info with password := literal "xxxxx" } else info) }
  (urlString changed).getD []

def urlRequestURI (receiver : Option URLRender) : Option URLBytes :=
  receiver.map fun u =>
    let main :=
      if u.opaquePart.isEmpty then
        let path := (escapedPath (some u.path)).getD []
        if path.isEmpty then [47] else path
      else if startsDoubleSlash u.opaquePart then u.scheme ++ [58] ++ u.opaquePart
      else u.opaquePart
    main ++ querySuffix u

structure URLAppend where
  output : URLBytes
  original : URLBytes
  reused : Bool

def urlAppendBinary (receiver : Option URLRender) (backing : URLBytes) (length : Nat) :
    Option URLAppend :=
  (urlString receiver).map fun content =>
    let output := backing.take length ++ content
    if output.length ≤ backing.length then
      ⟨output, output ++ backing.drop output.length, true⟩
    else ⟨output, backing, false⟩

def urlMarshalBinary (receiver : Option URLRender) : Option URLBytes :=
  (urlAppendBinary receiver [] 0).map URLAppend.output

theorem urlString_colon :
    urlString (some {path := ⟨[97, 58, 98], []⟩}) = some [46, 47, 97, 58, 98] := by decide

theorem urlRedacted_nil : urlRedacted none = [] := rfl

theorem urlRequestURI_empty : urlRequestURI (some {}) = some [47] := rfl

theorem urlAppendBinary_nil (backing : URLBytes) (length : Nat) :
    urlAppendBinary none backing length = none := rfl

theorem urlMarshalBinary_empty : urlMarshalBinary (some {}) = some [] := rfl

end GoLean.Net.Url
