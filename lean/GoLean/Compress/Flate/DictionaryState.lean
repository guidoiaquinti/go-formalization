import Init.Data.BitVec.Lemmas

namespace GoLean.Compress.Flate

structure DictionaryState where
  backing : List (BitVec 8)
  length : Nat
  read : Nat
  written : Nat
  full : Bool
  present : Bool

structure DictionaryView where
  bytes : List (BitVec 8)
  capacity : Nat
  present : Bool
  origin : Nat

private def view (s : DictionaryState) (start stop : Nat) : DictionaryView :=
  ⟨(s.backing.drop start).take (stop - start), s.backing.length - start, s.present, start⟩

def histSize (s : DictionaryState) : Nat := if s.full then s.length else s.written

def availRead (s : DictionaryState) : Nat := s.written - s.read

def availWrite (s : DictionaryState) : Nat := s.length - s.written

def writeSlice (s : DictionaryState) : DictionaryView := view s s.written s.length

def writeMark (s : DictionaryState) (count : Nat) : DictionaryState :=
  { s with written := s.written + count }

def writeByte (s : DictionaryState) (byte : BitVec 8) : Option DictionaryState :=
  if s.written < s.length then
    some { s with backing := s.backing.set s.written byte, written := s.written + 1 }
  else none

def readFlush (s : DictionaryState) : DictionaryView × DictionaryState :=
  (view s s.read s.written,
    if s.written = s.length then { s with read := 0, written := 0, full := true }
    else { s with read := s.written })

theorem histSize_full (s : DictionaryState) (full : s.full = true) :
    histSize s = s.length := by simp [histSize, full]

theorem availRead_ready (s : DictionaryState) (bounds : s.read ≤ s.written) :
    availRead s + s.read = s.written := Nat.sub_add_cancel bounds

theorem availWrite_ready (s : DictionaryState) (bounds : s.written ≤ s.length) :
    availWrite s + s.written = s.length := Nat.sub_add_cancel bounds

theorem writeSlice_origin (s : DictionaryState) : (writeSlice s).origin = s.written := rfl

theorem writeMark_value (s : DictionaryState) (n : Nat) :
    (writeMark s n).written = s.written + n := rfl

theorem writeByte_full (s : DictionaryState) (full : s.written = s.length) (b : BitVec 8) :
    writeByte s b = none := by simp [writeByte, full]

theorem readFlush_full (s : DictionaryState) (full : s.written = s.length) :
    (readFlush s).2 = { s with read := 0, written := 0, full := true } := by
  simp [readFlush, full]

end GoLean.Compress.Flate
