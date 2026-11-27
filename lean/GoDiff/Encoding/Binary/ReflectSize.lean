import GoLean.Encoding.Binary.ReflectSize
import GoDiff.Math.Bits.Tuple

open GoLean.Encoding.Binary

private def parseLayout (fuel : Nat) (tokens : List Int) : Option (Layout × List Int) :=
  match fuel, tokens with
  | 0, _ | _, [] => none
  | k+1, tag :: rest =>
    match tag with
    | 0 => some (.invalid, rest)
    | 1 => some (.unsupported, rest)
    | 2 => some (.scalar 1, rest)
    | 3 => some (.scalar 2, rest)
    | 4 => some (.scalar 4, rest)
    | 5 => some (.scalar 8, rest)
    | 6 => some (.scalar 16, rest)
    | 7 => some (.structEmpty, rest)
    | 8 => do
      let (field, rest) ← parseLayout k rest
      let (fields, tail) ← parseLayout k rest
      some (.structField field fields, tail)
    | 9 | 10 | 11 =>
      match rest with
      | [] => none
      | n :: rest => do
        let (child, tail) ← parseLayout k rest
        let layout := if tag == 9 then Layout.array n.toNat child
          else if tag == 10 then Layout.slice n.toNat child
          else Layout.pointer (n != 0) child
        some (layout, tail)
    | 12 | 13 => some (.scalar 1, rest)
    | 14 => some (.scalar 2, rest)
    | 15 | 16 | 18 => some (.scalar 4, rest)
    | 17 | 19 | 20 => some (.scalar 8, rest)
    | 21 | 22 => some (.unsupported, rest)
    | _ => none

private def wire (size : Option Nat) : String :=
  (size.map toString).getD "-1"

def main : IO Unit := runTuple fun row => do
  let some (layout, []) := parseLayout row.length row |
    throw (IO.userError "invalid layout")
  let byType := match sizeof layout with
    | .error _ => "panic"
    | .ok size => wire size
  return s!"{byType}:{wire (dataSize layout)}:{wire (binarySize layout)}"
