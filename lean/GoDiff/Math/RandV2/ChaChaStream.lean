import GoLean.Math.RandV2.ChaChaStream
import GoDiff.Math.Bits.Tuple

open GoLean.Math.RandV2

private def wire (bytes : List (BitVec 8)) := String.intercalate "," (bytes.map (toString ∘ BitVec.toNat))

private def draw : Nat → ChaChaStream → List (BitVec 64) × ChaChaStream
  | 0, state => ([], state)
  | n+1, state =>
    let (word, state) := chaCha8Uint64 state
    let (rest, after) := draw n state
    (word :: rest, after)

def main : IO Unit := runTuple fun row => do
  let drawCount :: first :: second :: front :: raw := row |
    throw (IO.userError "expected stream lengths, seeds and backing")
  let seed ← (raw.take 4).mapM (wordInput 64)
  let newSeed ← ((raw.drop 4).take 4).mapM (wordInput 64)
  let backing ← (raw.drop 8).mapM (wordInput 8)
  let initial := newChaCha8 seed
  let (values, s) := draw drawCount.toNat initial
  let (a, s) := chaCha8Read s first.toNat
  let (middle, s) := chaCha8Uint64 s
  let (b, s) := chaCha8Read s second.toNat
  let appended := chaCha8AppendBinary s front.toNat backing
  let reset := chaCha8Seed s newSeed
  let (resetWord, _) := chaCha8Uint64 reset
  return String.intercalate "|" [wire (chaCha8MarshalBinary initial),
    String.intercalate "," (values.map (toString ∘ BitVec.toNat)), wire a, toString middle.toNat, wire b,
    wire (chaCha8MarshalBinary s), s!"{appended.reused}:{wire appended.output}:{wire appended.original}",
    wire (chaCha8MarshalBinary reset), toString resetWord.toNat]
