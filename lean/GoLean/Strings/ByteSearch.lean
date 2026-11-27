import GoLean.Bytes.ByteSearch

namespace GoLean.Strings

/-- Go string search scans raw bytes.
    https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L99-L104
    https://github.com/golang/go/blob/go1.27.0/src/internal/stringslite/strings.go#L24-L26
    https://github.com/golang/go/blob/go1.27.0/src/internal/bytealg/indexbyte_generic.go#L22-L29 -/
def indexByte := Bytes.indexByte

/-- https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L265-L270
    https://github.com/golang/go/blob/go1.27.0/src/internal/bytealg/lastindexbyte_generic.go#L16-L23 -/
def lastIndexByte := Bytes.lastIndexByte

theorem indexByte_found (xs : List UInt8) (key : UInt8) (i : Nat) :
    indexByte xs key = some i ↔
      ∃ hi : i < xs.length, key = xs[i] ∧
        ∀ j (hj : j < i), key ≠ xs[j]'(Nat.lt_trans hj hi) :=
  Bytes.indexByte_found xs key i

theorem indexByte_missing (xs : List UInt8) (key : UInt8) :
    indexByte xs key = none ↔ key ∉ xs :=
  Bytes.indexByte_missing xs key

theorem lastIndexByte_found (xs : List UInt8) (key : UInt8) (i : Nat) :
    lastIndexByte xs key = some i ↔
      ∃ hi : i < xs.length, key = xs[i] ∧ ∀ j (hj : j < xs.length), i < j → key ≠ xs[j] :=
  Bytes.lastIndexByte_found xs key i

theorem lastIndexByte_missing (xs : List UInt8) (key : UInt8) :
    lastIndexByte xs key = none ↔ key ∉ xs :=
  Bytes.lastIndexByte_missing xs key

end GoLean.Strings
