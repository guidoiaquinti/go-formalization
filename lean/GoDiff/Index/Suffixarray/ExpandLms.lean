import GoLean.Index.Suffixarray.ExpandLms
import GoDiff.Math.Bits.Tuple

open GoLean.Index.Suffixarray

private def runs : List (BitVec width) → List (BitVec width × Nat)
  | [] => []
  | x :: xs =>
    match runs xs with
    | (y, n) :: rest => if x = y then (x, n + 1) :: rest else (x, 1) :: (y, n) :: rest
    | [] => [(x, 1)]

private def wire (values : List (BitVec width)) :=
  String.intercalate "," ((runs values).map fun (x, n) => s!"{x.toInt}.{n}")

private def suite
    (run : List Int → Nat → ExpandState width → ExpandState width)
    (bins mode : Nat) (text : List Int) (compacted : List Int) : String := Id.run do
  let tallies : Array (BitVec width) := text.foldl
      (fun bins c => bins.set! c.toNat (bins[c.toNat]! + 1)) (Array.replicate bins 0)
  let f := if mode == 0 then none else some
    ((if mode == 1 then (-1 : BitVec width) :: List.replicate (bins - 1) 17 else tallies.toList) ++ [93,94])
  let b := List.replicate bins 19 ++ [91,92]
  let sa := ((compacted.map (BitVec.ofInt width)) ++ List.replicate (text.length - compacted.length) (-7)).toArray
  let result := run text compacted.length ⟨⟨⟨f, if mode == 0 then 0 else bins⟩, ⟨some b, bins⟩⟩, sa⟩
  return wire result.suffixes.toList ++ "/" ++ (result.buckets.frequency.backing.map wire).getD "nil" ++
    "/" ++ (result.buckets.bucket.backing.map wire).getD "nil"

def main : IO Unit := runTuple fun input => do
  let mode :: bins :: size :: rest := input | throw (IO.userError "expected mode, bins and text size")
  let byteText := rest.take size.toNat
  let wideText := byteText.map (· * 2)
  let compacted := rest.drop size.toNat
  return suite expandByte32 256 mode.toNat byteText compacted ++ "|" ++
    suite expandByte64 256 mode.toNat byteText compacted ++ "|" ++
    suite expand32 bins.toNat mode.toNat wideText compacted ++ "|" ++
    suite expand64 bins.toNat mode.toNat wideText compacted
