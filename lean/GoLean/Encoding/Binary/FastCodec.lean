import GoLean.Encoding.Binary.EndianRead
import GoLean.Encoding.Binary.EndianPut
import GoLean.Encoding.Binary.IntDataSize

namespace GoLean.Encoding.Binary

inductive FastForm
  | scalar | pointer | slice | unsupported
  deriving DecidableEq

private def boolBit (value : Nat) : Nat := if value = 0 then 0 else 1

private def decodeCell (big : Bool) (kind : FastKind) (cells : List (BitVec 8))
    (index : Nat) : Option Nat := do
  let width := fastWidth kind
  let raw ← readWord big width (cells.drop (index * width))
  pure (if kind = .bool then boolBit raw else raw)

private def decodedPrefix (big : Bool) (kind : FastKind) (cells : List (BitVec 8))
    (index : Nat) : List Nat → Bool × List Nat
  | [] => (false, [])
  | value :: rest =>
    match decodeCell big kind cells index with
    | none => (true, value :: rest)
    | some value =>
      let (panic, next) := decodedPrefix big kind cells (index + 1) rest
      (panic, value :: next)

/-- Status is unsupported, success, or bounds panic. Values are raw words. -/
def decodeFast (big : Bool) (kind : FastKind) (form : FastForm)
    (cells : List (BitVec 8)) (values : List Nat) : String × List Nat :=
  if kind = .unsupported ∨ form = .unsupported ∨ form = .scalar then ("unsupported", values)
  else
    let count :=
      if form = .pointer then 1
      else if kind = .uint8 then min values.length cells.length
      else if kind = .bool ∨ kind = .int8 then min values.length cells.length
      else values.length
    let (failed, decoded) := decodedPrefix big kind cells 0 (values.take count)
    let byteLengthPanic := form = .slice ∧ (kind = .bool ∨ kind = .int8) ∧ values.length < cells.length
    (if failed ∨ byteLengthPanic then "bounds" else "ok", decoded ++ values.drop count)

private def encodedWords (big : Bool) (kind : FastKind) :
    List Nat → List (BitVec 8) → Bool × List (BitVec 8)
  | [], cells => (false, cells)
  | value :: rest, cells =>
    let width := fastWidth kind
    let word := if kind = .bool then boolBit value else value
    let (ok, written) := putWord big width word cells
    if !ok then (true, cells)
    else
      let (failed, after) := encodedWords big kind rest (written.drop width)
      (failed, written.take width ++ after)

def encodeFast (big : Bool) (kind : FastKind) (form : FastForm)
    (values : List Nat) (cells : List (BitVec 8)) : String × List (BitVec 8) :=
  if kind = .unsupported ∨ form = .unsupported then ("ok", cells)
  else
    let selected := if form = .slice ∧ kind = .uint8 then values.take cells.length else values
    let (failed, after) := encodedWords big kind selected cells
    (if failed then "bounds" else "ok", after)

theorem decodeFast_value_is_unsupported (big : Bool) (kind : FastKind)
    (bytes : List (BitVec 8)) (values : List Nat) :
    decodeFast big kind .scalar bytes values = ("unsupported", values) := by
  simp [decodeFast]

theorem encodeFast_unsupported (big : Bool) (kind : FastKind)
    (values : List Nat) (bytes : List (BitVec 8)) :
    encodeFast big kind .unsupported values bytes = ("ok", bytes) := by
  simp [encodeFast]

end GoLean.Encoding.Binary
