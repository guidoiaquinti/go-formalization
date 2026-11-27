import GoLean.Index.Suffixarray.Accessors

namespace GoLean.Index.Suffixarray

structure BucketState (width : Nat) where
  frequency : View (BitVec width)
  bucket : View (BitVec width)

private def cells (view : View α) := view.backing.getD []

private def save (view : View α) (backing : List α) : View α :=
  { view with backing := view.backing.map (fun _ => backing) }

private def countSymbols (length : Nat) (text : List Int) (counts : List (BitVec width)) :
    Bool × List (BitVec width) :=
  match text with
  | [] => (true, counts)
  | c :: rest =>
    if c < 0 || length ≤ c.toNat then (false, counts)
    else
      match counts[c.toNat]? with
      | none => (false, counts)
      | some old => countSymbols length rest (counts.set c.toNat (old + 1))

private def frequency (byteInput : Bool) (text : List Int) (state : BucketState width) :
    Option (View (BitVec width)) × BucketState width :=
  let separate := state.frequency.backing.isSome
  let original := if separate then state.frequency else state.bucket
  if separate && original.length = 0 then (none, state)
  else if separate && 0 ≤ ((cells original).headD 0).toInt then (some original, state)
  else
    let length := if byteInput then 256 else original.length
    let backing := cells original
    if backing.length < length then (none, state)
    else
      let initial := List.replicate length 0 ++ backing.drop length
      let (ok, counted) := countSymbols length text initial
      let updated := save original counted
      let state := if separate then { state with frequency := updated }
        else { state with bucket := updated }
      (if ok then some { updated with length := length } else none, state)

private def prefixTotals (rightEdge : Bool) (total : BitVec width) :
    List (BitVec width) → List (BitVec width)
  | [] => []
  | n :: rest =>
    let next := total + n
    (if rightEdge then next else total) :: prefixTotals rightEdge next rest

private def bucketize (byteInput rightEdge : Bool) (text : List Int) (state : BucketState width) :
    Bool × BucketState width :=
  let (counts, state) := frequency byteInput text state
  match counts with
  | none => (false, state)
  | some view =>
    let length := if byteInput then 256 else view.length
    let backing := cells state.bucket
    let target := if byteInput then 256 else state.bucket.length
    if byteInput && ((cells view).length < length || backing.length < target) then (false, state)
    else
      let output := (prefixTotals rightEdge 0 ((cells view).take length)).take target
      let data := output ++ backing.drop output.length
      (length ≤ target, { state with bucket := save state.bucket data })

/-- Go 1.27.0 src/index/suffixarray/sais.go:213–228. Byte data uses 256 bins. -/
def freqByte32 (text : List Int) (state : BucketState 32) := frequency true text state

/-- Go 1.27.0 src/index/suffixarray/sais2.go:202–216. Byte data uses 64-bit bins. -/
def freqByte64 (text : List Int) (state : BucketState 64) := frequency true text state

/-- Go 1.27.0 src/index/suffixarray/sais2.go:218–231. The selected length is preserved. -/
def freq32 (text : List Int) (state : BucketState 32) := frequency false text state

/-- Go 1.27.0 src/index/suffixarray/sais2.go:233–246. The selected length is preserved. -/
def freq64 (text : List Int) (state : BucketState 64) := frequency false text state

/-- Go 1.27.0 src/index/suffixarray/sais.go:231–240. Emit before the addition. -/
def bucketMinByte32 (text : List Int) (state : BucketState 32) := bucketize true false text state

/-- Go 1.27.0 src/index/suffixarray/sais2.go:248–257. Emit before the addition. -/
def bucketMinByte64 (text : List Int) (state : BucketState 64) := bucketize true false text state

/-- Go 1.27.0 src/index/suffixarray/sais2.go:259–266. Emit before the addition. -/
def bucketMin32 (text : List Int) (state : BucketState 32) := bucketize false false text state

/-- Go 1.27.0 src/index/suffixarray/sais2.go:268–275. Emit before the addition. -/
def bucketMin64 (text : List Int) (state : BucketState 64) := bucketize false false text state

/-- Go 1.27.0 src/index/suffixarray/sais.go:246–255. Emit after the addition. -/
def bucketMaxByte32 (text : List Int) (state : BucketState 32) := bucketize true true text state

/-- Go 1.27.0 src/index/suffixarray/sais2.go:277–286. Emit after the addition. -/
def bucketMaxByte64 (text : List Int) (state : BucketState 64) := bucketize true true text state

/-- Go 1.27.0 src/index/suffixarray/sais2.go:288–295. Emit after the addition. -/
def bucketMax32 (text : List Int) (state : BucketState 32) := bucketize false true text state

/-- Go 1.27.0 src/index/suffixarray/sais2.go:297–304. Emit after the addition. -/
def bucketMax64 (text : List Int) (state : BucketState 64) := bucketize false true text state

private def cached (width n : Nat) : BucketState width :=
  ⟨⟨some (List.replicate n 3), n⟩, ⟨some (List.replicate n 9), n⟩⟩

theorem freqByte32_cached : freqByte32 [255] (cached 32 256) =
    (some (cached 32 256).frequency, cached 32 256) := rfl
theorem freqByte64_cached : freqByte64 [255] (cached 64 256) =
    (some (cached 64 256).frequency, cached 64 256) := rfl
theorem freq32_cached : freq32 [1] (cached 32 2) =
    (some (cached 32 2).frequency, cached 32 2) := rfl
theorem freq64_cached : freq64 [1] (cached 64 2) =
    (some (cached 64 2).frequency, cached 64 2) := rfl

set_option maxRecDepth 4096 in
theorem bucketMinByte32_first : ((bucketMinByte32 [] (cached 32 256)).2.bucket.backing.map List.head?) =
    some (some 0) := rfl
set_option maxRecDepth 4096 in
theorem bucketMinByte64_first : ((bucketMinByte64 [] (cached 64 256)).2.bucket.backing.map List.head?) =
    some (some 0) := rfl
theorem bucketMin32_prefix : (bucketMin32 [] (cached 32 2)).2.bucket.backing = some [0, 3] := rfl
theorem bucketMin64_prefix : (bucketMin64 [] (cached 64 2)).2.bucket.backing = some [0, 3] := rfl

set_option maxRecDepth 4096 in
theorem bucketMaxByte32_first : ((bucketMaxByte32 [] (cached 32 256)).2.bucket.backing.map List.head?) =
    some (some 3) := rfl
set_option maxRecDepth 4096 in
theorem bucketMaxByte64_first : ((bucketMaxByte64 [] (cached 64 256)).2.bucket.backing.map List.head?) =
    some (some 3) := rfl
theorem bucketMax32_prefix : (bucketMax32 [] (cached 32 2)).2.bucket.backing = some [3, 6] := rfl
theorem bucketMax64_prefix : (bucketMax64 [] (cached 64 2)).2.bucket.backing = some [3, 6] := rfl

end GoLean.Index.Suffixarray
