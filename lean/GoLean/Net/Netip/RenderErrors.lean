import GoLean.Net.Netip.ParseFields
import GoLean.Strconv.PrintRanges
import GoLean.Unicode.Utf8.DecodeRuneInString

namespace GoLean.Net.Netip

private def printable (r : Nat) : Bool :=
  let rec search (fuel lo hi : Nat) : Bool :=
    match fuel with
    | 0 => false
    | fuel + 1 =>
      if lo ≥ hi then false else
        let mid := (lo + hi) / 2
        let (first, last) := Strconv.printRanges[mid]!
        if r < first then search fuel lo mid
        else if r > last then search fuel (mid + 1) hi
        else true
  search (Nat.log2 Strconv.printRanges.size + 1) 0 Strconv.printRanges.size

private def escapedHex (kind value width : Nat) : ParseBytes :=
  92 :: kind :: ((List.range width).reverse.map fun i =>
    let digit := value / 16^i % 16
    if digit < 10 then 48 + digit else 87 + digit)

private def quoteRune (r : Nat) (bytes : ParseBytes) : ParseBytes :=
  if r == 34 || r == 92 then [92, r]
  else if printable r then bytes
  else match r with
    | 7 => [92, 97]
    | 8 => [92, 98]
    | 9 => [92, 116]
    | 10 => [92, 110]
    | 11 => [92, 118]
    | 12 => [92, 102]
    | 13 => [92, 114]
    | _ =>
      if r < 32 || r == 127 then escapedHex 120 r 2
      else if r < 65536 then escapedHex 117 r 4
      else escapedHex 85 r 8

private def quoteBody : Nat → ParseBytes → ParseBytes → ParseBytes
  | 0, _, acc => acc
  | _ + 1, [], acc => acc
  | fuel + 1, lead :: rest, acc =>
    let (r, width) := Unicode.Utf8.decodeRuneInString (lead :: rest)
    let next := if r == 65533 && width == 1 then escapedHex 120 lead 2
      else quoteRune r ((lead :: rest).take width)
    quoteBody fuel ((lead :: rest).drop width) (next.reverse ++ acc)

private def quote (input : ParseBytes) : ParseBytes :=
  34 :: (34 :: quoteBody input.length input []).reverse

private def ascii (s : String) : ParseBytes := s.toUTF8.data.toList.map UInt8.toNat

def renderAddrError (input message suffix : ParseBytes) : ParseBytes :=
  ascii "ParseAddr(" ++ quote input ++ ascii "): " ++ message ++
    if suffix.isEmpty then [] else ascii " (at " ++ quote suffix ++ [41]

def renderPrefixError (input message : ParseBytes) : ParseBytes :=
  ascii "netip.ParsePrefix(" ++ quote input ++ ascii "): " ++ message

theorem renderAddrError_empty :
    renderAddrError [] [] [] = ascii "ParseAddr(\"\"): " := by decide

theorem renderPrefixError_empty :
    renderPrefixError [] [] = ascii "netip.ParsePrefix(\"\"): " := by decide

end GoLean.Net.Netip
