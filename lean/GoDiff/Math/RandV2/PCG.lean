import GoLean.Math.RandV2.PCG
import GoDiff.Math.Bits.Tuple

open GoLean.Math.RandV2

private def wire (bytes : List (BitVec 8)) :=
  String.intercalate "," (bytes.map (toString ∘ BitVec.toNat))

private def stateWire (s : PCGState) := s!"{s.hi.toNat},{s.lo.toNat}"

def main : IO Unit := runTuple fun row => do
  let hi :: lo :: seedHi :: seedLo :: length :: capacity :: raw := row |
    throw (IO.userError "expected state, seeds, geometry and unmarshal input")
  let s := newPCG (← wordInput 64 hi) (← wordInput 64 lo)
  let seeded := pcgSeed s (← wordInput 64 seedHi) (← wordInput 64 seedLo)
  let bytes ← raw.mapM (wordInput 8)
  let next := pcgNext s
  let (first, after) := pcgUint64 s
  let (second, afterAgain) := pcgUint64 after
  let (ok, restored) := pcgUnmarshalBinary s (bytes.drop capacity.toNat)
  let out := pcgAppendBinary s length.toNat (bytes.take capacity.toNat)
  return String.intercalate "|" [
    stateWire s, stateWire seeded, stateWire next, s!"{first.toNat},{second.toNat}:{stateWire afterAgain}",
    wire (pcgMarshalBinary s), s!"{out.reused}:{wire out.output}:{wire out.original}",
    s!"{ok}:{stateWire restored}"]
