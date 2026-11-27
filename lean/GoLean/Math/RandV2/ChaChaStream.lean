import GoLean.Encoding.Binary.EndianPut

namespace GoLean.Math.RandV2

open GoLean.Encoding.Binary

private def rotate (x : BitVec 32) (distance : Nat) : BitVec 32 :=
  (x <<< distance) ||| (x >>> (32 - distance))

private def quarter (s : Array (BitVec 32)) (i j k l : Nat) : Array (BitVec 32) :=
  let a := s[i]! + s[j]!
  let d := rotate (s[l]! ^^^ a) 16
  let c := s[k]! + d
  let b := rotate (s[j]! ^^^ c) 12
  let a := a + b
  let d := rotate (d ^^^ a) 8
  let c := c + d
  let b := rotate (b ^^^ c) 7
  (((s.set! i a).set! j b).set! k c).set! l d

private def doubleRound (s : Array (BitVec 32)) : Array (BitVec 32) :=
  let columns := [(0,4,8,12), (1,5,9,13), (2,6,10,14), (3,7,11,15)]
  let diagonals := [(0,5,10,15), (1,6,11,12), (2,7,8,13), (3,4,9,14)]
  (columns ++ diagonals).foldl (fun state (i,j,k,l) => quarter state i j k l) s

private def block (seed : List (BitVec 64)) (counter : Nat) : Array (BitVec 64) :=
  let key := seed.flatMap fun w => [BitVec.ofNat 32 w.toNat, BitVec.ofNat 32 (w.toNat / 2^32)]
  let blocks := (List.range 4).map fun lane =>
    let initial := ([0x61707865, 0x3320646e, 0x79622d32, 0x6b206574] ++ key ++
      [BitVec.ofNat 32 (counter + lane), 0, 0, 0]).toArray
    let mixed := (List.replicate 4 ()).foldl (fun s _ => doubleRound s) initial
    (List.range 16).map fun i => if 4 ≤ i && i < 12 then mixed[i]! + initial[i]! else mixed[i]!
  let words := (List.range 16).flatMap fun i => blocks.map fun lanes => lanes[i]!
  ((List.range 32).map fun i =>
    BitVec.ofNat 64 (words[2*i]!.toNat + 2^32 * words[2*i+1]!.toNat)).toArray

structure ChaChaStream where
  seed : List (BitVec 64)
  generated : Array (BitVec 64)
  counter : Nat
  index : Nat
  pending : List (BitVec 8)

def newChaCha8 (seed : List (BitVec 64)) : ChaChaStream :=
  ⟨seed, block seed 0, 0, 0, []⟩

def chaCha8Seed (_ : ChaChaStream) (seed : List (BitVec 64)) := newChaCha8 seed

private def refill (state : ChaChaStream) : ChaChaStream :=
  let reseed := state.counter = 12
  let seed := if reseed then (state.generated.toList).drop 28 else state.seed
  let counter := if reseed then 0 else state.counter + 4
  { state with seed, counter, index := 0, generated := block seed counter }

def chaCha8Uint64 (state : ChaChaStream) : BitVec 64 × ChaChaStream :=
  let limit := if state.counter = 12 then 28 else 32
  let state := if state.index ≥ limit then refill state else state
  (state.generated[state.index]!, {state with index := state.index + 1})

private def bytesOf (big : Bool) (word : BitVec 64) : List (BitVec 8) :=
  (putWord big 8 word.toNat (List.replicate 8 0)).2

private def fill : Nat → Nat → ChaChaStream → List (BitVec 8) × ChaChaStream
  | 0, _, state => ([], state)
  | fuel+1, count, state =>
    if count = 0 then ([], state) else
    let (word, after) := chaCha8Uint64 state
    let bytes := bytesOf false word
    if count < 8 then (bytes.take count, {after with pending := bytes.drop count})
    else
      let (rest, done) := fill fuel (count - 8) after
      (bytes ++ rest, done)

def chaCha8Read (state : ChaChaStream) (count : Nat) : List (BitVec 8) × ChaChaStream :=
  let used := state.pending.take count
  let (rest, after) := fill (count / 8 + 1) (count - used.length)
    {state with pending := state.pending.drop count}
  (used ++ rest, after)

private def stateBytes (s : ChaChaStream) : List (BitVec 8) :=
  [99,104,97,99,104,97,56,58] ++ bytesOf true (BitVec.ofNat 64 (s.counter / 4 * 32 + s.index)) ++
    s.seed.flatMap (bytesOf false)

private def appendChunks (s : ChaChaStream) : List (List (BitVec 8)) :=
  let header := if s.pending.isEmpty then []
    else [[114,101,97,100,98,117,102,58], [BitVec.ofNat 8 s.pending.length], s.pending]
  header ++ [stateBytes s]

def chaCha8MarshalBinary (s : ChaChaStream) : List (BitVec 8) := (appendChunks s).flatten

structure ChaChaAppendResult where
  output : List (BitVec 8)
  original : List (BitVec 8)
  reused : Bool

private def fittingChunks (spare : Nat) : List (List (BitVec 8)) → List (BitVec 8)
  | [] => []
  | first :: rest =>
    if first.length ≤ spare then first ++ fittingChunks (spare - first.length) rest else []

def chaCha8AppendBinary (s : ChaChaStream) (length : Nat) (cells : List (BitVec 8)) :
    ChaChaAppendResult :=
  let output := chaCha8MarshalBinary s
  let written := fittingChunks (cells.length - length) (appendChunks s)
  ⟨cells.take length ++ output, cells.take length ++ written ++ cells.drop (length + written.length),
    length + output.length ≤ cells.length⟩

theorem newChaCha8_pending (seed : List (BitVec 64)) : (newChaCha8 seed).pending = [] := rfl

theorem chaCha8Seed_reset (old : ChaChaStream) (seed : List (BitVec 64)) :
    chaCha8Seed old seed = newChaCha8 seed := rfl

theorem chaCha8Uint64_pending (s : ChaChaStream) :
    (chaCha8Uint64 s).2.pending = s.pending := by
  by_cases h : (if s.counter = 12 then 28 else 32) ≤ s.index
  · simp [chaCha8Uint64, h, refill]
  · simp [chaCha8Uint64, h]

theorem chaCha8Read_zero (s : ChaChaStream) : (chaCha8Read s 0).1 = [] := rfl

theorem chaCha8MarshalBinary_unbuffered (s : ChaChaStream) (empty : s.pending = []) :
    chaCha8MarshalBinary s = stateBytes s := by
  simp [chaCha8MarshalBinary, appendChunks, empty]

theorem chaCha8AppendBinary_output (s : ChaChaStream) (length : Nat) (b : List (BitVec 8)) :
    (chaCha8AppendBinary s length b).output = b.take length ++ chaCha8MarshalBinary s := rfl

end GoLean.Math.RandV2
