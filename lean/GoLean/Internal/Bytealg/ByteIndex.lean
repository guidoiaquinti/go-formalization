import GoLean.Bytes.ByteSearch

namespace GoLean.Internal.Bytealg

/-- Go 1.27.0 src/internal/bytealg/indexbyte_native.go:10. None represents -1. -/
def indexByte (bytes : List UInt8) (target : UInt8) : Option Nat :=
  Bytes.indexByte bytes target

theorem indexByte_found (bytes : List UInt8) (target : UInt8) (i : Nat) :
    indexByte bytes target = some i ↔
      ∃ hi : i < bytes.length, target = bytes[i] ∧ ∀ j (hj : j < i), target ≠ bytes[j]'(Nat.lt_trans hj hi) :=
  Bytes.indexByte_found bytes target i

theorem indexByte_missing (bytes : List UInt8) (target : UInt8) :
    indexByte bytes target = none ↔ target ∉ bytes :=
  Bytes.indexByte_missing bytes target

/-- Go 1.27.0 src/internal/bytealg/indexbyte_native.go:13. None represents -1. -/
def indexByteString (bytes : List UInt8) (target : UInt8) : Option Nat :=
  Bytes.indexByte bytes target

theorem indexByteString_found (bytes : List UInt8) (target : UInt8) (i : Nat) :
    indexByteString bytes target = some i ↔
      ∃ hi : i < bytes.length, target = bytes[i] ∧ ∀ j (hj : j < i), target ≠ bytes[j]'(Nat.lt_trans hj hi) :=
  Bytes.indexByte_found bytes target i

theorem indexByteString_missing (bytes : List UInt8) (target : UInt8) :
    indexByteString bytes target = none ↔ target ∉ bytes :=
  Bytes.indexByte_missing bytes target

/-- Go 1.27.0 src/internal/bytealg/lastindexbyte_generic.go:7–14. None represents -1. -/
def lastIndexByte (bytes : List UInt8) (target : UInt8) : Option Nat :=
  Bytes.lastIndexByte bytes target

theorem lastIndexByte_found (bytes : List UInt8) (target : UInt8) (i : Nat) :
    lastIndexByte bytes target = some i ↔
      ∃ hi : i < bytes.length, target = bytes[i] ∧ ∀ j (hj : j < bytes.length), i < j → target ≠ bytes[j] :=
  Bytes.lastIndexByte_found bytes target i

theorem lastIndexByte_missing (bytes : List UInt8) (target : UInt8) :
    lastIndexByte bytes target = none ↔ target ∉ bytes :=
  Bytes.lastIndexByte_missing bytes target

/-- Go 1.27.0 src/internal/bytealg/lastindexbyte_generic.go:16–23. None represents -1. -/
def lastIndexByteString (bytes : List UInt8) (target : UInt8) : Option Nat :=
  Bytes.lastIndexByte bytes target

theorem lastIndexByteString_found (bytes : List UInt8) (target : UInt8) (i : Nat) :
    lastIndexByteString bytes target = some i ↔
      ∃ hi : i < bytes.length, target = bytes[i] ∧ ∀ j (hj : j < bytes.length), i < j → target ≠ bytes[j] :=
  Bytes.lastIndexByte_found bytes target i

theorem lastIndexByteString_missing (bytes : List UInt8) (target : UInt8) :
    lastIndexByteString bytes target = none ↔ target ∉ bytes :=
  Bytes.lastIndexByte_missing bytes target

end GoLean.Internal.Bytealg
