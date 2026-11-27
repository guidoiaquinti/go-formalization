import GoLean.Index.Suffixarray.Accessors
import GoLean.Encoding.Binary.VarintPut
import GoLean.Encoding.Binary.VarintDecode

namespace GoLean.Index.Suffixarray

open GoLean.Encoding.Binary (putVarint putUvarint varint)

structure IntegerWrite where
  status : String
  count : Nat
  backing : List (BitVec 8)
  offered : Option (Nat × Nat × List (BitVec 8))

private def writeOffer (capacity count : Nat) (writer : String) (backing : List (BitVec 8))
    (size : Nat) : IntegerWrite :=
  if backing.length < size then ⟨"bounds", 0, backing, none⟩
  else if writer = "nilIO" then ⟨"nilIO", 0, backing, none⟩
  else ⟨writer, count, backing, some (size, capacity, backing.take size)⟩

/-- Go 1.27.0 src/index/suffixarray/suffixarray.go:89–93. Write the whole
    ten-byte header, including retained padding, and return the writer error. -/
def writeInt (backing : List (BitVec 8)) (length : Nat) (value : BitVec 64)
    (writer : String) : IntegerWrite :=
  let (written, updated) := putVarint value length backing
  if written.isNone then ⟨"bounds", 0, updated, none⟩
  else writeOffer backing.length 0 writer updated 10

structure IntegerRead where
  status : String
  value : Int
  backing : List (BitVec 8)
  requests : List Nat

private def readRequests (readLimit : Nat) : Nat → Nat → Nat → List Nat
  | 0, _, _ => []
  | fuel + 1, needed, available =>
    if needed = 0 then []
    else
      let read := min (min needed readLimit) available
      needed :: (if read = available then []
        else readRequests readLimit fuel (needed - read) (available - read))

/-- Go 1.27.0 src/index/suffixarray/suffixarray.go:96–100. ReadFull writes
    through a capacity-checked ten-byte view, then Varint reads the original
    visible buffer even when ReadFull returned an error. -/
def readInt (backing : List (BitVec 8)) (length : Nat) (source : List (BitVec 8))
    (readLimit : Nat) (terminal : String) : IntegerRead :=
  if backing.length < 10 then ⟨"bounds", 0, backing, []⟩
  else if terminal = "nilIO" then ⟨"nilIO", 0, backing, []⟩
  else
    let count := min 10 source.length
    let updated := source.take count ++ backing.drop count
    let status := if count = 10 then "ok"
      else if terminal = "eof" && count ≠ 0 then "unexpected" else terminal
    ⟨status, (varint (updated.take length)).1, updated,
      readRequests (max 1 readLimit) 11 10 source.length⟩

private def packSlice (data : Ints) (length : Nat) :
    Nat → Nat → Nat → List (BitVec 8) → Option (Nat × Nat) × List (BitVec 8)
  | 0, count, pos, backing => (some (count, pos), backing)
  | fuel + 1, count, pos, backing =>
    if length < pos + 10 then (some (count, pos), backing)
    else
      match intsGet (some data) (Int.ofNat count) with
      | .error _ => (none, backing)
      | .ok value =>
        let (written, output) := putUvarint (BitVec.ofInt 64 value) (length - pos) (backing.drop pos)
        let backing := backing.take pos ++ output
        match written with
        | none => (none, backing)
        | some n => packSlice data length fuel (count + 1) (pos + n) backing

/-- Go 1.27.0 src/index/suffixarray/suffixarray.go:104–120. Reserve ten
    bytes for the size and require room for a full word before each entry. -/
def writeSlice (backing : List (BitVec 8)) (length : Nat) (data : Ints)
    (writer : String) : IntegerWrite :=
  let entries := (intsLen (some data)).toOption.getD 0
  let (packed, updated) := packSlice data length entries 0 10 backing
  match packed with
  | none => ⟨"bounds", 0, updated, none⟩
  | some (count, size) =>
    let (head, final) := putVarint (BitVec.ofNat 64 size) length updated
    if head.isNone then ⟨"bounds", 0, final, none⟩
    else writeOffer backing.length count writer final size

theorem writeInt_keepsPadding :
    writeInt (List.replicate 12 99) 1 (BitVec.ofInt 64 (-1)) "other" =
      ⟨"other", 0, 1 :: List.replicate 11 99, some (10, 12, 1 :: List.replicate 9 99)⟩ := by
  simp [writeInt, putVarint, GoLean.Encoding.Binary.signedVarintWord,
    GoLean.Encoding.Binary.uvarintBytes]
  rfl

theorem readInt_staleAfterError :
    readInt (List.replicate 10 1) 10 [] 3 "eof" =
      ⟨"eof", -1, List.replicate 10 1, [10]⟩ := rfl

theorem readInt_shortView :
    readInt (List.replicate 10 99) 0 [14] 1 "eof" =
      ⟨"unexpected", 0, 14 :: List.replicate 9 99, [10]⟩ := rfl

theorem writeSlice_wordRoom :
    let data : Ints := ⟨⟨some [1], 1⟩, ⟨none, 0⟩⟩
    (writeSlice (List.replicate 20 99) 19 data "ok").count = 0 := by
  let header := putVarint (BitVec.ofNat 64 10) 19 (List.replicate 20 99)
  change (if header.1.isNone then IntegerWrite.mk "bounds" 0 header.2 none
    else writeOffer 20 0 "ok" header.2 10).count = 0
  split
  · rfl
  · unfold writeOffer
    split <;> rfl

end GoLean.Index.Suffixarray
