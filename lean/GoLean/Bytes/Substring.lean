import GoLean.Bytes.Affix
import GoLean.Strings.LastIndex

namespace GoLean.Bytes

private def raw (buffer : Slice) : List Nat := (contents buffer).map BitVec.toNat

/-- The model scans candidate byte windows in ascending order. Go optimizes
    the candidates with length guards, byte search and rolling hashes.
    https://github.com/golang/go/blob/go1.27.0/src/bytes/bytes.go#L1287-L1367 -/
def index (buffer sep : Slice) : Option Nat := Strings.index (raw buffer) (raw sep)

/-- Go's backward search chooses the greatest matching byte offset.
    https://github.com/golang/go/blob/go1.27.0/src/bytes/bytes.go#L102-L120 -/
def lastIndex (buffer sep : Slice) : Option Nat := Strings.lastIndex (raw buffer) (raw sep)

/-- https://github.com/golang/go/blob/go1.27.0/src/bytes/bytes.go#L76-L79 -/
def contains (buffer sep : Slice) : Bool := (index buffer sep).isSome

private def splitAtMatch (buffer : Slice) (size : Nat) : Option Nat → Slice × Slice × Bool
  | none => (buffer, none, false)
  | some i =>
    (buffer.map fun (data, capacity) => (data.take i, capacity),
      buffer.map fun (data, capacity) => (data.drop (i + size), capacity - (i + size)), true)

/-- https://github.com/golang/go/blob/go1.27.0/src/bytes/bytes.go#L1376-L1386 -/
def cut (buffer sep : Slice) : Slice × Slice × Bool :=
  splitAtMatch buffer (contents sep).length (index buffer sep)

/-- https://github.com/golang/go/blob/go1.27.0/src/bytes/bytes.go#L1424-L1435 -/
def cutLast (buffer sep : Slice) : Slice × Slice × Bool :=
  splitAtMatch buffer (contents sep).length (lastIndex buffer sep)

theorem index_spec (buffer sep : Slice) :
    Strings.FirstMatch (raw buffer) (raw sep) (index buffer sep) :=
  Strings.index_spec _ _

theorem lastIndex_spec (buffer sep : Slice) :
    Strings.LastMatch (raw buffer) (raw sep) (lastIndex buffer sep) :=
  Strings.lastIndex_spec _ _

theorem contains_spec (buffer sep : Slice) :
    contains buffer sep = true ↔ ∃ i, Strings.matchesAt (raw buffer) (raw sep) i :=
  Strings.contains_spec _ _

theorem cut_spec (buffer sep : Slice) :
    match index buffer sep with
    | none => cut buffer sep = (buffer, none, false)
    | some i =>
      cut buffer sep =
        (buffer.map fun (data, capacity) => (data.take i, capacity),
          buffer.map fun (data, capacity) =>
            (data.drop (i + (contents sep).length), capacity - (i + (contents sep).length)), true) := by
  cases h : index buffer sep <;> simp [cut, h, splitAtMatch]

theorem cutLast_spec (buffer sep : Slice) :
    match lastIndex buffer sep with
    | none => cutLast buffer sep = (buffer, none, false)
    | some i =>
      cutLast buffer sep =
        (buffer.map fun (data, capacity) => (data.take i, capacity),
          buffer.map fun (data, capacity) =>
            (data.drop (i + (contents sep).length), capacity - (i + (contents sep).length)), true) := by
  cases h : lastIndex buffer sep <;> simp [cutLast, h, splitAtMatch]

end GoLean.Bytes
