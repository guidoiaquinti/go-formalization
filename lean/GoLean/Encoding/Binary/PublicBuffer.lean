import GoLean.Encoding.Binary.FastCodec
import GoLean.Encoding.Binary.Ensure

namespace GoLean.Encoding.Binary

structure BufferResult (α : Type) where
  status : String
  count : Nat
  state : α

private def basicSize (kind : FastKind) (form : FastForm) (words : List Nat) : Nat :=
  match form with
  | .scalar | .pointer => fastWidth kind
  | .slice => words.length * fastWidth kind
  | .unsupported => 0

private def badType (kind : FastKind) (form : FastForm) : Bool :=
  kind == .unsupported || form == .unsupported

def binaryDecode (big namedByte : Bool) (kind : FastKind) (form : FastForm)
    (cells : List (BitVec 8)) (words : List Nat) : BufferResult (List Nat) :=
  let n := basicSize kind form words
  if badType kind form || (namedByte && form == .scalar) then ⟨"type", 0, words⟩
  else if cells.length < n then ⟨"short", 0, words⟩
  else if form = .scalar then ⟨"type", 0, words⟩
  else if n = 0 then ⟨"ok", 0, words⟩
  else
    let (status, values) := decodeFast big kind form cells words
    ⟨status, if status == "ok" then n else 0, values⟩

def binaryEncode (big : Bool) (kind : FastKind) (form : FastForm)
    (words : List Nat) (cells : List (BitVec 8)) : BufferResult (List (BitVec 8)) :=
  let n := basicSize kind form words
  if badType kind form then ⟨"type", 0, cells⟩
  else if cells.length < n then ⟨"short", 0, cells⟩
  else
    let (status, values) := encodeFast big kind form words cells
    ⟨status, n, values⟩

structure BufferAppend where
  status : String
  output : List (BitVec 8)
  backing : List (BitVec 8)
  allocated : Bool
  nilOutput : Bool

def binaryAppend (big : Bool) (kind : FastKind) (form : FastForm) (words : List Nat)
    (view : GoLean.Slices.GrowView (BitVec 8)) : BufferAppend :=
  if badType kind form then ⟨"type", [], view.cells, false, true⟩
  else
    let n := basicSize kind form words
    match ensure view (Int.ofNat n) with
    | .error _ => ⟨"bounds", [], view.cells, false, true⟩
    | .ok grown =>
      let (status, payload) := encodeFast big kind form words grown.position
      let out := grown.buffer.take view.length ++ payload
      let backing := if grown.allocated then view.cells else out ++ view.cells.drop out.length
      ⟨status, out, backing, grown.allocated, grown.nilBuffer⟩

theorem binaryDecode_short (big namedByte : Bool) (kind : FastKind) (form : FastForm)
    (cells : List (BitVec 8)) (words : List Nat)
    (known : (badType kind form || (namedByte && form == .scalar)) = false)
    (short : cells.length < basicSize kind form words) :
    (binaryDecode big namedByte kind form cells words).state = words := by
  simp [binaryDecode, known, short]

theorem binaryEncode_short (big : Bool) (kind : FastKind) (form : FastForm)
    (cells : List (BitVec 8)) (words : List Nat)
    (known : badType kind form = false) (short : cells.length < basicSize kind form words) :
    (binaryEncode big kind form words cells).state = cells := by
  simp [binaryEncode, known, short]

theorem binaryAppend_bad_type (big : Bool) (kind : FastKind) (words : List Nat)
    (view : GoLean.Slices.GrowView (BitVec 8)) :
    (binaryAppend big kind .unsupported words view).backing = view.cells := by
  simp [binaryAppend, badType]

end GoLean.Encoding.Binary
