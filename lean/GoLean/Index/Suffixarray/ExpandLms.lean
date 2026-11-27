import GoLean.Index.Suffixarray.Buckets

namespace GoLean.Index.Suffixarray

structure ExpandState (width : Nat) where
  buckets : BucketState width
  suffixes : Array (BitVec width)

private def reserve (text : Array Int) (sa : Array (BitVec width))
    (bucket : Array (BitVec width)) (index : Nat) :
    BitVec width × BitVec width × Array (BitVec width) :=
  let slot := sa[index]!
  let symbol := text[slot.toNat]!.toNat
  let target := bucket[symbol]! - 1
  (slot, target, bucket.set! symbol target)

private def sweep (text : Array Int) :
    Nat → Nat → BitVec width → BitVec width → Array (BitVec width) →
      Array (BitVec width) → Array (BitVec width) × Array (BitVec width)
  | 0, _, _, _, bucket, sa => (bucket, sa)
  | i + 1, source, slot, target, bucket, sa =>
    if i == target.toNat then
      let sa := sa.set! i slot
      if source == 0 then sweep text i source slot target bucket sa
      else
        let (nextSlot, nextTarget, bucket) := reserve text sa bucket (source - 1)
        sweep text i (source - 1) nextSlot nextTarget bucket sa
    else sweep text i source slot target bucket (sa.set! i 0)

private def expand (upper : List Int → BucketState width → Bool × BucketState width)
    (text : List Int) (num : Nat) (s : ExpandState width) : ExpandState width :=
  let (_, buckets) := upper text s.buckets
  let bucket := (buckets.bucket.backing.getD []).toArray
  let (slot, target, bucket) := reserve text.toArray s.suffixes bucket (num - 1)
  let (bucket, sa) := sweep text.toArray s.suffixes.size (num - 1) slot target bucket s.suffixes
  ⟨{buckets with bucket := {buckets.bucket with backing := some bucket.toList}}, sa⟩

def expandByte32 := expand bucketMaxByte32
def expandByte64 := expand bucketMaxByte64
def expand32 := expand bucketMax32
def expand64 := expand bucketMax64

private def exampleState (width bins : Nat) : ExpandState width :=
  ⟨⟨⟨none, 0⟩, ⟨some (List.replicate bins 9), bins⟩⟩, #[1, 9, 9]⟩

set_option maxRecDepth 4096 in
theorem expandByte32_places :
    (expandByte32 [2, 1, 3] 1 (exampleState 32 256)).suffixes = #[1, 0, 0] := rfl

set_option maxRecDepth 4096 in
theorem expandByte64_places :
    (expandByte64 [2, 1, 3] 1 (exampleState 64 256)).suffixes = #[1, 0, 0] := rfl

theorem expand32_places :
    (expand32 [2, 1, 3] 1 (exampleState 32 4)).suffixes = #[1, 0, 0] := rfl

theorem expand64_places :
    (expand64 [2, 1, 3] 1 (exampleState 64 4)).suffixes = #[1, 0, 0] := rfl

end GoLean.Index.Suffixarray
