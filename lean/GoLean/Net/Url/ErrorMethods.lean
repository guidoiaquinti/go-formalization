import GoLean.Net.Url.ErrorQuotePrint
import GoLean.Unicode.Utf8.DecodeRuneInString

namespace GoLean.Net.Url

abbrev ErrorBytes := List Nat

private def printSearch (r : Nat) : Nat → Nat → Nat → Bool
  | 0, _, _ => false
  | fuel + 1, lo, hi =>
    if hi ≤ lo then false else
      let mid := (lo + hi) / 2
      let (first, last) := goQuotePrintRanges[mid]!
      if r < first then printSearch r fuel lo mid
      else if r > last then printSearch r fuel (mid + 1) hi
      else true

private def quotePrint (r : Nat) : Bool :=
  printSearch r 12 0 goQuotePrintRanges.size

private def hexDigits (r : Nat) : Nat → ErrorBytes
  | 0 => []
  | n + 1 =>
    let v := (r / 16^n) % 16
    (if v < 10 then 48 + v else 87 + v) :: hexDigits r n

private def quoteRune (r : Nat) (original : ErrorBytes) : ErrorBytes :=
  if r == 34 || r == 92 then [92, r]
  else if quotePrint r then original
  else match r with
    | 7 => [92, 97]
    | 8 => [92, 98]
    | 9 => [92, 116]
    | 10 => [92, 110]
    | 11 => [92, 118]
    | 12 => [92, 102]
    | 13 => [92, 114]
    | n =>
      if n < 32 || n == 127 then [92, 120] ++ hexDigits n 2
      else if n < 65536 then [92, 117] ++ hexDigits n 4
      else [92, 85] ++ hexDigits n 8

private def quoteBody : Nat → ErrorBytes → ErrorBytes
  | 0, _ => []
  | _ + 1, [] => []
  | fuel + 1, bytes@(b :: _) =>
    let (r, width) := GoLean.Unicode.Utf8.decodeRuneInString bytes
    let head := if width == 1 && r == 65533 then [92, 120] ++ hexDigits b 2
      else quoteRune r (bytes.take width)
    head ++ quoteBody fuel (bytes.drop width)

private def quoted (bytes : ErrorBytes) : ErrorBytes :=
  34 :: quoteBody bytes.length bytes ++ [34]

private def ascii (s : String) : ErrorBytes := s.toList.map Char.toNat

inductive ErrorFault where
  | nilReceiver
  | callback (token : Nat)

structure ErrorPayload where
  innerIdentity : Nat
  operation : ErrorBytes
  url : ErrorBytes
  formattedError : ErrorBytes
  timeout : Option (Except Nat Bool)
  temporary : Option (Except Nat Bool)

def urlErrorError (receiver : Option ErrorPayload) : Except ErrorFault ErrorBytes :=
  match receiver with
  | none => .error .nilReceiver
  | some e => .ok (e.operation ++ [32] ++ quoted e.url ++ [58, 32] ++ e.formattedError)

def urlErrorUnwrap (receiver : Option ErrorPayload) : Except ErrorFault Nat :=
  match receiver with
  | none => .error .nilReceiver
  | some e => .ok e.innerIdentity

private def dispatch (receiver : Option ErrorPayload)
    (method : ErrorPayload → Option (Except Nat Bool)) : Except ErrorFault Bool :=
  match receiver with
  | none => .error .nilReceiver
  | some e =>
    match method e with
    | none => .ok false
    | some (.ok result) => .ok result
    | some (.error token) => .error (.callback token)

def urlErrorTimeout (e : Option ErrorPayload) := dispatch e ErrorPayload.timeout

def urlErrorTemporary (e : Option ErrorPayload) := dispatch e ErrorPayload.temporary

def escapeErrorError (bytes : ErrorBytes) : ErrorBytes := ascii "invalid URL escape " ++ quoted bytes

def invalidHostErrorError (bytes : ErrorBytes) : ErrorBytes :=
  ascii "invalid character " ++ quoted bytes ++ ascii " in host name"

theorem urlErrorError_empty :
    urlErrorError (some ⟨0, [], [], [], none, none⟩) = .ok [32, 34, 34, 58, 32] := rfl

theorem urlErrorUnwrap_identity (e : ErrorPayload) : urlErrorUnwrap (some e) = .ok e.innerIdentity := rfl

theorem urlErrorTimeout_missing (e : ErrorPayload) (h : e.timeout = none) :
    urlErrorTimeout (some e) = .ok false := by simp [urlErrorTimeout, dispatch, h]

theorem urlErrorTemporary_missing (e : ErrorPayload) (h : e.temporary = none) :
    urlErrorTemporary (some e) = .ok false := by simp [urlErrorTemporary, dispatch, h]

theorem escapeErrorError_empty : escapeErrorError [] = ascii "invalid URL escape \"\"" := rfl

theorem invalidHostErrorError_empty :
    invalidHostErrorError [] = ascii "invalid character \"\" in host name" := rfl

end GoLean.Net.Url
