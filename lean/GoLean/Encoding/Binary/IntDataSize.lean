namespace GoLean.Encoding.Binary

inductive FastKind
  | bool | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
  | float32 | float64 | unsupported
  deriving DecidableEq

def fastWidth : FastKind → Nat
  | .bool | .int8 | .uint8 => 1
  | .int16 | .uint16 => 2
  | .int32 | .uint32 | .float32 => 4
  | .int64 | .uint64 | .float64 => 8
  | .unsupported => 0

structure ByteView where
  backing : List (BitVec 8)
  length : Nat
  nilSlice : Bool

inductive FastData
  | atom (kind : FastKind)
  | pointer (kind : FastKind)
  | elements (kind : FastKind) (length : Nat)
  | bytes (view : ByteView)
  | unsupported

def intDataSize : FastData → Nat × Option ByteView
  | .atom kind | .pointer kind => (fastWidth kind, none)
  | .elements kind length => (fastWidth kind * length, none)
  | .bytes bytes => (bytes.length, some bytes)
  | .unsupported => (0, none)

theorem intDataSize_byte_alias (view : ByteView) :
    intDataSize (.bytes view) = (view.length, some view) := rfl

theorem intDataSize_pointer (kind : FastKind) :
    intDataSize (.pointer kind) = intDataSize (.atom kind) := rfl

theorem intDataSize_unsupported : intDataSize .unsupported = (0, none) := rfl

end GoLean.Encoding.Binary
