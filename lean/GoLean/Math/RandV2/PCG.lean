import GoLean.Encoding.Binary.EndianRead
import GoLean.Encoding.Binary.EndianPut

namespace GoLean.Math.RandV2

open GoLean.Encoding.Binary

structure PCGState where
  hi : BitVec 64
  lo : BitVec 64

def newPCG (hi lo : BitVec 64) : PCGState := ⟨hi, lo⟩

def pcgSeed (_ : PCGState) (hi lo : BitVec 64) : PCGState := newPCG hi lo

private def join (s : PCGState) : Nat := s.hi.toNat * 2^64 + s.lo.toNat

def pcgNext (s : PCGState) : PCGState :=
  let mult := 2549297995355413924 * 2^64 + 4865540595714422341
  let inc := 6364136223846793005 * 2^64 + 1442695040888963407
  let raw := BitVec.ofNat 128 (join s * mult + inc)
  ⟨BitVec.ofNat 64 (raw.toNat / 2^64), BitVec.ofNat 64 raw.toNat⟩

def pcgUint64 (s : PCGState) : BitVec 64 × PCGState :=
  let next := pcgNext s
  let hi := (next.hi ^^^ (next.hi >>> 32)) * 0xda942042e4dd58b5
  (((hi ^^^ (hi >>> 48)) * (next.lo ||| 1)), next)

private def magic : List (BitVec 8) := [112, 99, 103, 58]

private def bytesOf (word : BitVec 64) : List (BitVec 8) :=
  (bigPutUint64 word (List.replicate 8 0)).2

def pcgMarshalBinary (state : PCGState) : List (BitVec 8) :=
  magic ++ bytesOf state.hi ++ bytesOf state.lo

structure PCGAppendResult where
  output : List (BitVec 8)
  original : List (BitVec 8)
  reused : Bool

def pcgAppendBinary (state : PCGState) (length : Nat) (backing : List (BitVec 8)) :
    PCGAppendResult :=
  let encoded := pcgMarshalBinary state
  let remaining := backing.length - length
  let written := if remaining ≥ 20 then 20 else if remaining ≥ 12 then 12
    else if remaining ≥ 4 then 4 else 0
  let initial := backing.take length
  ⟨initial ++ encoded, initial ++ encoded.take written ++ backing.drop (length + written),
    remaining ≥ 20⟩

def pcgUnmarshalBinary (state : PCGState) (input : List (BitVec 8)) : Bool × PCGState :=
  if input.length = 20 && input.take 4 == magic then
    let hi := BitVec.ofNat 64 ((bigUint64 (input.drop 4)).getD 0)
    let lo := BitVec.ofNat 64 ((bigUint64 (input.drop 12)).getD 0)
    (true, ⟨hi, lo⟩)
  else (false, state)

theorem newPCG_seed (hi lo : BitVec 64) :
    (newPCG hi lo).hi = hi ∧ (newPCG hi lo).lo = lo := ⟨rfl, rfl⟩

theorem pcgSeed_reset (old : PCGState) (hi lo : BitVec 64) :
    pcgSeed old hi lo = newPCG hi lo := rfl

theorem pcgNext_zero : pcgNext (newPCG 0 0) = newPCG 6364136223846793005 1442695040888963407 :=
  rfl

theorem pcgUint64_updates (state : PCGState) : (pcgUint64 state).2 = pcgNext state := rfl

theorem pcgMarshalBinary_magic (state : PCGState) :
    (pcgMarshalBinary state).take 4 = magic := by simp [pcgMarshalBinary, magic]

theorem pcgAppendBinary_output (state : PCGState) (length : Nat) (cells : List (BitVec 8)) :
    (pcgAppendBinary state length cells).output = cells.take length ++ pcgMarshalBinary state := rfl

theorem pcgUnmarshalBinary_bad_length (state : PCGState) (bytes : List (BitVec 8))
    (bad : bytes.length ≠ 20) : pcgUnmarshalBinary state bytes = (false, state) := by
  simp [pcgUnmarshalBinary, bad]

end GoLean.Math.RandV2
