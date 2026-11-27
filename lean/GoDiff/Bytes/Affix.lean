import GoLean.Bytes.Affix
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def describe : Slice → String
  | none => "nil"
  | some (bytes, cap) => String.intercalate " " (toString cap :: bytes.map (fun b => toString b.toNat))

def main : IO Unit := runTuple fun fields => do
  let cap :: n :: words := fields | throw (IO.userError "expected capacity, affix length and bytes")
  if cap < -1 || n < 0 || words.length < n.toNat then
    throw (IO.userError "invalid affix input")
  let bytes ← words.mapM (wordInput 8)
  let affix := bytes.take n.toNat
  let data := bytes.drop n.toNat
  if (cap == -1 && !data.isEmpty) || (cap ≥ 0 && cap.toNat < data.length) then
    throw (IO.userError "invalid slice capacity")
  let src : Slice := if cap == -1 then none else some (data, cap.toNat)
  let (leading, leadingFound) := cutPrefix src affix
  let (trailing, trailingFound) := cutSuffix src affix
  return String.intercalate "|" [
    toString (hasPrefix src affix), toString (hasSuffix src affix),
    toString leadingFound ++ " " ++ describe leading, toString trailingFound ++ " " ++ describe trailing,
    describe (trimPrefix src affix), describe (trimSuffix src affix)]
