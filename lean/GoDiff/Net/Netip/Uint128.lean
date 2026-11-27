import GoLean.Net.Netip.Uint128
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Netip

private def wire (u : Uint128) : String := s!"{u.hi.toNat},{u.lo.toNat}"
private def optWire (u : Option Uint128) : String := (u.map wire).getD "panic"

def main : IO Unit := runTuple fun row => do
  let [n, bit, missing, hi, lo, vh, vl] := row |
    throw (IO.userError "expected mask count, bit index, receiver and two words")
  let bit ← wordInput 8 bit
  let u : Uint128 := ⟨← wordInput 64 hi, ← wordInput 64 lo⟩
  let v : Uint128 := ⟨← wordInput 64 vh, ← wordInput 64 vl⟩
  let halves := match uint128Halves (if missing = 0 then some u else none) v with
    | none => "panic"
    | some (old, new) => wire old ++ ":" ++ wire new
  return String.intercalate "|" [optWire (mask6 n), toString (uint128IsZero u),
    wire (uint128And u v), wire (uint128Xor u v), wire (uint128Or u v), wire (uint128Not u),
    wire (uint128SubOne u), wire (uint128AddOne u), halves,
    optWire (uint128BitsSetFrom u bit), optWire (uint128BitsClearedFrom u bit)]
