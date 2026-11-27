import GoLean.Encoding.Binary.EndianRead
import GoLean.Encoding.Binary.EndianPut

namespace GoLean.Encoding.Binary

inductive ValueItem where
  | field (width components : Nat) (boolean : Bool)
  | blank (width : Nat)
  | ignored

private def available (width length offset capacity : Nat) : Bool :=
  if width = 1 then offset < length else offset + width ≤ capacity

private def putCell (big : Bool) (width offset word : Nat)
    (cells : List (BitVec 8)) : List (BitVec 8) :=
  cells.take offset ++ (putWord big width word (cells.drop offset)).2

private def readComponents (big : Bool) (width length : Nat) (cells : List (BitVec 8)) :
    Nat → Nat → Option (List Nat) × Nat
  | 0, offset => (some [], offset)
  | count + 1, offset =>
    if available width length offset cells.length then
      let value := (readWord big width (cells.drop offset)).getD 0
      let (rest, after) := readComponents big width length cells count (offset + width)
      (rest.map (value :: ·), after)
    else (none, offset)

private def writeComponents (big : Bool) (width length : Nat) :
    List Nat → Nat → List (BitVec 8) → Bool × Nat × List (BitVec 8)
  | [], offset, cells => (true, offset, cells)
  | word :: rest, offset, cells =>
    if available width length offset cells.length then
      writeComponents big width length rest (offset + width) (putCell big width offset word cells)
    else (false, offset, cells)

def decoderValue (big : Bool) (length : Nat) (cells : List (BitVec 8)) :
    List ValueItem → Nat → List Nat → Bool × Nat × List Nat
  | [], offset, words => (true, offset, words)
  | .ignored :: rest, offset, words => decoderValue big length cells rest offset words
  | .blank width :: rest, offset, words => decoderValue big length cells rest (offset + width) words
  | .field width count boolean :: rest, offset, words =>
    let (field, next) := readComponents big width length cells count offset
    match field with
    | none => (false, next, words)
    | some field =>
      let value := if boolean then field.map (fun n => if n = 0 then 0 else 1) else field
      let (ok, after, tail) := decoderValue big length cells rest next (words.drop count)
      (ok, after, value ++ tail)

def encoderValue (big : Bool) (length : Nat) :
    List ValueItem → Nat → List Nat → List (BitVec 8) → Bool × Nat × List (BitVec 8)
  | [], offset, _, cells => (true, offset, cells)
  | .ignored :: rest, offset, words, cells => encoderValue big length rest offset words cells
  | .blank width :: rest, offset, words, cells =>
    if offset + width ≤ cells.length then
      let cleared := cells.take offset ++ List.replicate width 0 ++ cells.drop (offset + width)
      encoderValue big length rest (offset + width) words cleared
    else (false, offset, cells)
  | .field width count boolean :: rest, offset, words, cells =>
    let input := words.take count
    let input := if boolean then input.map (fun n => if n = 0 then 0 else 1) else input
    let (ok, next, after) := writeComponents big width length input offset cells
    if ok then encoderValue big length rest next (words.drop count) after
    else (false, next, after)

theorem decoderValue_blank (big : Bool) (length width offset : Nat) (cells : List (BitVec 8))
    (rest : List ValueItem) (words : List Nat) :
    decoderValue big length cells (.blank width :: rest) offset words =
      decoderValue big length cells rest (offset + width) words :=
  rfl

theorem encoderValue_blank (big : Bool) (length width offset : Nat) (cells : List (BitVec 8))
    (words : List Nat) (fits : offset + width ≤ cells.length) :
    encoderValue big length [.blank width] offset words cells =
      (true, offset + width, cells.take offset ++ List.replicate width 0 ++ cells.drop (offset + width)) := by
  simp [encoderValue, fits]

end GoLean.Encoding.Binary
