import GoLean.Math.Big.IntWords
import GoDiff.Math.Bits.Tuple

open GoLean.Math.Big

private def wire (n : IntWords width) :=
  s!"{n.negative}:{String.intercalate "," (n.words.map (toString ∘ BitVec.toNat))}"

private def renderOption [ToString α] (x : Option α) := x.map toString |>.getD "panic"

def main : IO Unit := runTuple fun fields => do
  let width :: present :: negative :: raw :: words := fields |
    throw (IO.userError "expected native width, receiver and sign, constructor word and magnitude")
  let width := width.toNat
  let live := present != 0
  let x := if live then some (⟨negative != 0, words.map (BitVec.ofInt width)⟩ : IntWords width) else none
  let raw := BitVec.ofInt 64 raw
  return String.intercalate "|" [wire (newInt width raw),
    ((setInt64 width live raw).map wire).getD "panic", ((setUint64 width live raw).map wire).getD "panic",
    renderOption (intSign x), renderOption (intToInt64 x), renderOption (intToUint64 x), renderOption (intIsInt64 x),
    renderOption (intIsUint64 x), renderOption (intBitLen x)]
